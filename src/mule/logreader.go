package mule

import (
	//"compute"
	"bufio"
	"github.com/anaray/compute"
	"github.com/anaray/regnet"
	"io"
	"net/http"
	"os"
	"time"
)

type LogReaderCompute struct {
	//Regnet *regnet.Regnet
}

type Logs struct {
	Store string
}

var logReaderLogger *compute.Log

func LogReader() *LogReaderCompute {
	//r, _ := regnet.New()
	//return &LogReaderCompute{Regnet: r}
	return &LogReaderCompute{}
}

func (reader *LogReaderCompute) String() string { return "compute.logreader" }

func (reader *LogReaderCompute) Execute(arg *compute.Args) {
	reg, _ := regnet.New()
	arg.Container["regnet"] = reg
	logReaderLogger = arg.Logger
	err := reg.AddPatternsFromFile("/home/msi/Desktop/metricstream.regnet")
	if err != nil {
		logReaderLogger.Logf("ERROR:", err)
		os.Exit(1)
	}

	logReaderLogger.Logf("Creating LogReader HTTP listener at port %s", "8080")
	handler := func(w http.ResponseWriter, r *http.Request) {
		filePath := r.FormValue("path")
		go reader.process(filePath, arg)
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func (reader *LogReaderCompute) process(file string, arg *compute.Args) { //, regexp *regexp.Regexp) {
	logReaderLogger.Logf("Received request to parse file %s", file)
	//initialize http handler and listen for POST message, get file name from request

	defer func() {
		if r := recover(); r != nil {
			logReaderLogger.Logf("ERROR: %v", r)
		}
	}()

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	// make a read buffer
	r := bufio.NewReaderSize(f, 1024*1024)
	var lg *Logs
	logReaderLogger.Logf("Started parsing file %s at %s:", file, time.Now().String())
	for {
		// read a chunk
		line, err := r.ReadSlice('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			packet := compute.NewPacket()
			packet["log"] = lg
			packet["source"] = file
			arg.Outgoing <- packet
			logReaderLogger.Logf("Completed parsing file %s at %s:", file, time.Now().String())
			break
		}
		//exists, _ := reader.Regnet.Exists(line, "%{MS_DELIM}")
		regnet := arg.Container["regnet"].(*regnet.Regnet)
		exists, _ := regnet.Exists(line, "%{MS_DELIM}")
		if exists {
			if lg != nil && len(lg.Store) > 0 {
				//push it further
				packet := compute.NewPacket()
				packet["log"] = lg
				packet["source"] = file
				arg.Outgoing <- packet
			}
			lg = NewLog(line)
		} else {
			if len(lg.Store) > 0 {
				store := []byte(lg.Store)
				lg.Store = string(append(store[:], line[:]...))
			}
		}
	}
}

func NewLog(line []byte) *Logs {
	log := new(Logs)
	log.Store = string(line)
	return log
}
