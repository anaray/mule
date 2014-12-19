package compute

import (
	"bufio"
	"github.com/anaray/regnet"
	"io"
	"net/http"
	"os"
	"time"
)

type LogReaderCompute struct {
	Regnet *regnet.Regnet
}

type Logs struct {
	Store string
}

var logger *Log

func LogReader() *LogReaderCompute {
	r, _ := regnet.New()
	r.AddPattern("MS_DATE_TIME", `((\d\d[\/]*){3}[\s]*(\d\d[:]*){3}.(\d*))`)
	r.AddPattern("MS_TZ", `((?:[PMCEI][SD]T|UTC)|GMT-\d\d:\d\d)`)
	r.AddPattern("LOGLEVEL", `([Aa]lert|ALERT|[Tt]race|TRACE|[Dd]ebug|DEBUG|[Nn]otice|NOTICE|[Ii]nfo|INFO|[Ww]arn?(?:ing)?|WARN?(?:ING)?|[Ee]rr?(?:or)?|ERR?(?:OR)?|[Cc]rit?(?:ical)?|CRIT?(?:ICAL)?|[Ff]atal|FATAL|[Ss]evere|SEVERE|EMERG(?:ENCY)?|[Ee]merg(?:ency)?)`)
	r.AddPattern("MS_DELIM", `%{MS_DATE_TIME}\s%{MS_TZ}\s%{LOGLEVEL}`)
	return &LogReaderCompute{Regnet: r}
}

func (reader *LogReaderCompute) String() string { return "compute.LogReaderCompute" }

func (reader *LogReaderCompute) Execute(arg Args) {
	logger = arg.Logger
	logger.logf("Creating LogReader HTTP listener at port %s","8080")
	handler := func(w http.ResponseWriter, r *http.Request) {
		filePath := r.FormValue("path")
		go reader.process(filePath, arg)
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func (reader *LogReaderCompute) process(file string, arg Args) { //, regexp *regexp.Regexp) {
	logger.logf("Recevied request to parse file %s ::", file)
	//initialize http handler and listen for POST message, get file name from request

	defer func() {
		if r := recover(); r != nil {
			logger.logf("ERROR: %v", r)
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
	logger.logf("Started parsing file %s at %s:", file, time.Now().String())
	for {
		// read a chunk
		line, err := r.ReadSlice('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			packet := NewPacket()
			packet["log"] = lg
			packet["source"] = file
			arg.Outgoing <- packet
			logger.logf("Completed parsing file %s at %s:", file, time.Now().String())
			break
		}

		exists, _ := reader.Regnet.Exists(line, "%{MS_DELIM}")
		if exists {
			if lg != nil && len(lg.Store) > 0 {
				//push it further
				packet := NewPacket()
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
