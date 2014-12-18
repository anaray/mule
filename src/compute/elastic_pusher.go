package compute

import (
	"bytes"
	"encoding/json"
	"fmt"
	elastigo "github.com/mattbaird/elastigo/lib"
	"os"
	"strconv"
	"time"
	"container/list"
)

type ElasticSearchPlugin struct {}

func ElasticSearch() *ElasticSearchPlugin {
	return new(ElasticSearchPlugin)
}

type PacketHolder struct {
	buffer *list.List
}

func (p *ElasticSearchPlugin) Execute(arg Args) {
	packetChannel := make(chan Packet, 10000)
	connection := elastigo.NewConn()
	connection.Domain = "localhost"
	indexer := connection.NewBulkIndexer(10)
	indexer.BulkMaxDocs = 10000
	indexer.BufferDelayMax = 1000 * time.Millisecond
	indexer.Sender = func(buf *bytes.Buffer) error {
		return indexer.Send(buf)
	}
	defer indexer.Stop()
	indexer.Start()

	go pushToElasticSearch(packetChannel, indexer)
	for {
		packet := <-arg.Incoming
		packetChannel<-packet
	}
}

func pushToElasticSearch(packetChannel chan Packet, indexer *elastigo.BulkIndexer) {
	for {
		select {
		case packet := <-packetChannel:
			log := packet["log"].(*Logs)	
			src := packet["source"]
			if log != nil {
				logRecord := LogRecord{Record: log.Store, Source: src.(string)}
				logRecordStr, _ := json.Marshal(logRecord)
				indexer.Index("testindex", "user", strconv.FormatInt(time.Now().UnixNano(), 36), "", nil, logRecordStr, false)
			}
		}
	}
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