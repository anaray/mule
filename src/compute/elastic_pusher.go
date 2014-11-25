package compute

import (
	"bytes"
	"encoding/json"
	"fmt"
	elastigo "github.com/mattbaird/elastigo/lib"
	"os"
	"strconv"
	"sync"
	"time"
	"strings"
)

type ElasticSearchPlugin struct {
}

func ElasticSearch() *ElasticSearchPlugin {
	return new(ElasticSearchPlugin)
}

type PacketHolder struct {
	buffer []Packet
}

func NewPacketHolder() PacketHolder { return PacketHolder{buffer: []Packet{}} }

func (holder PacketHolder) Merge(other PacketHolder) PacketHolder {
	return PacketHolder{buffer: append(holder.buffer, other.buffer...)}
}

func (p *ElasticSearchPlugin) Execute(arg Args) {
	in := make(chan PacketHolder, 10000)
	out := make(chan PacketHolder, 10000)
	go pushToElasticSearch(out)
	go coalesce(in, out)
	for {
		// packet is a Map
		packet := <-arg.Incoming
		in <- PacketHolder{buffer: []Packet{packet}}
	}
}

func coalesce(in chan PacketHolder, out chan PacketHolder) {
	packetHolder := NewPacketHolder()
	timer := time.NewTimer(0)

	var timerCh <-chan time.Time
	var outCh chan<- PacketHolder

	for {
		select {
		case ph := <-in:
			packetHolder = packetHolder.Merge(ph)
			if timerCh == nil {
				timer.Reset(500 * time.Millisecond)
				timerCh = timer.C
			}
		case <-timerCh:
			outCh = out
			timerCh = nil
		case outCh <- packetHolder:
			packetHolder = NewPacketHolder()
			outCh = nil
		}
	}
}

func pushToElasticSearch(out chan PacketHolder) {
	for {
		select {
		case holder := <-out:
			go process(holder)
		}
	}
}

var count int
func process(packetHolder PacketHolder) {
	var wg sync.WaitGroup
	connection := elastigo.NewConn()
	connection.Domain = "localhost"

	indexer := connection.NewBulkIndexer(10)
	indexer.BulkMaxDocs = 10000
	indexer.BufferDelayMax = 1000 * time.Millisecond
	indexer.Sender = func(buf *bytes.Buffer) error {
		wg.Done()
		return indexer.Send(buf)
	}
	indexer.Start()
	defer indexer.Stop()
	buffer := packetHolder.buffer
	for _, packet := range buffer {
		log := packet["log"].(Logs)
		src := packet["source"]
		//n:=bytes.Index(log.Store,[]byte{0})
		l:=string(log.Store[0:len(log.Store)])
		logRecord := LogRecord{Record: l, Source: src.(string)}
		logRecordStr, _ := json.Marshal(logRecord)
		wg.Add(1)
		indexer.Index("testindex", "user", strconv.FormatInt(time.Now().UnixNano(), 36), "", nil, logRecordStr, false)
	}
	wg.Wait()
}

type LogRecord struct {
	Record string `json:"record"`
	Source string `json:"source"`
}

func exitIfErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		//os.Exit(1)
	}
}
