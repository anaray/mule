package compute

import (
	"bytes"
	"encoding/json"
	"fmt"
	elastigo "github.com/mattbaird/elastigo/lib"
	"os"
	"strconv"
	"time"
)

var logger *Log

type ElasticSearchCompute struct{}

func ElasticSearch() *ElasticSearchCompute {
	return new(ElasticSearchCompute)
}

func (e *ElasticSearchCompute) String() string { return "compute.ElasticSearchCompute" }

func (e *ElasticSearchCompute) Execute(arg Args) {
	logger = arg.Logger
	packetChannel := make(chan Packet, 10000)
	connection := elastigo.NewConn()
	connection.Domain = "localhost"
	//indexer := connection.NewBulkIndexer(10)
	indexer := connection.NewBulkIndexerErrors(10, 60)
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
		packetChannel <- packet
	}
}

func pushToElasticSearch(packetChannel chan Packet, indexer *elastigo.BulkIndexer) {
	//loop the error channel of the elastigo bulk indexer to check if there are any
	// errors
	go func() {
  		for errBuf := range indexer.ErrorChannel {
    		logger.logf("ERROR: %v",errBuf.Err)
  		}
	}()

	for {
		select {
		case packet := <-packetChannel:
			log := packet["log"].(*Logs)
			src := packet["source"]
			if log != nil {
				logRecord := LogRecord{Record: log.Store, Source: src.(string)}
				logRecordStr, _ := json.Marshal(logRecord)
				err := indexer.Index("testindex", "user", strconv.FormatInt(time.Now().UnixNano(), 36), "", nil, logRecordStr, false)
				if err != nil {
					logger.logf("ERROR: %v",err)
				}
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
