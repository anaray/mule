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
	buffer := holder.buffer
	fmt.Println("inside Merge..")
	for _, packet := range buffer {
		log := packet["log"].(Logs)
		//fmt.Print(">>>>>>>>>>>>>>>>> ", string(log.Store))
		//src := packet["source"]
		//n:=bytes.Index(log.Store,[]byte{0})
		l:=string(log.Store[0:len(log.Store)])
		if !strings.HasPrefix(l,"07/09/14"){
			count=count+1
			fmt.Println("Jumbled:", l)
			
			fmt.Println("Jumbled count ::",count)

			fmt.Println("Actual")
			for _, p:=range log.Store {
				fmt.Print(string(p))
			}
			fmt.Println("===========================================================")
		}
	}
	p := PacketHolder{buffer: append(holder.buffer, other.buffer...)}
	//return PacketHolder{buffer: append(holder.buffer, other.buffer...)}

	/*buf := p.buffer
	for _, packet := range buf {
		log := packet["log"].(Logs)
		//fmt.Print(">>>>>>>>>>>>>>>>> ", string(log.Store))
		//src := packet["source"]
		//n:=bytes.Index(log.Store,[]byte{0})
		l:=string(log.Store[0:len(log.Store)])
		if !strings.HasPrefix(l,"07/09/14"){
			count=count+1
			fmt.Println("Jumbled:", string(log.Store))
			
			fmt.Println("Jumbled count ::",count)

			/*fmt.Println("Actual")
			for _, p:=range log.Store {
				fmt.Print(string(p))
			}
			fmt.Println("===========================================================")
		}
		
	
	}*/


	return p
}

func (p *ElasticSearchPlugin) Execute(arg Args) {
	in := make(chan PacketHolder, 10000)
	out := make(chan PacketHolder, 10000)
	go pushToElasticSearch(out)
	go coalesce(in, out)
	for {
		// packet is a Map
		packet := <-arg.Incoming
				//packet["log"] = lg
				//packet["source"] = file
				//fmt.Println("checking packet .....",string(lg.Store[0:8]) )
				lg:=packet["log"].(Logs)
				if  string(lg.Store[0:8])!="07/09/14" {
					fmt.Println("FOUND ::", string(lg.Store[0:8]))
				}else {
					//fmt.Println("Not found ::", string(lg.Store[0:8]))
				}

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
			buffer := holder.buffer
			for _, packet := range buffer {
				log := packet["log"].(Logs)

				if !strings.HasPrefix(string(log.Store[0:len(log.Store)]),"07/09/14"){
					//fmt.Println("Jumbled >",string(log.Store))
				}
				//fmt.Println("::: >",string(log.Store))

			}
			//go process(holder)
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
		//fmt.Print(">>>>>>>>>>>>>>>>> ", string(log.Store))
		src := packet["source"]
		//n:=bytes.Index(log.Store,[]byte{0})
		l:=string(log.Store[0:len(log.Store)])
		if !strings.HasPrefix(l,"07/09/14"){
			count=count+1
			fmt.Println("Jumbled:", l)
			
			fmt.Println("Jumbled count ::",count)

			fmt.Println("Actual")
			for _, p:=range log.Store {
				fmt.Print(string(p))
			}
			fmt.Println("===========================================================")
		}
		
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
