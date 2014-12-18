package compute

import (
	"bufio"
	"io"
	"os"
	"fmt"
	"github.com/anaray/regnet"
	"net/http"
	"regexp"
	"time"
	"strings"
	"sync"
)

type LogReaderCompute struct {
	Regnet *regnet.Regnet
}

type Logs struct {
	Store string
}

func LogReader() *LogReaderCompute {
	r, _ := regnet.New()
	r.AddPattern("MS_DATE_TIME", `((\d\d[\/]*){3}[\s]*(\d\d[:]*){3}.(\d*))`)
	r.AddPattern("MS_TZ", `((?:[PMCEI][SD]T|UTC)|GMT-\d\d:\d\d)`)
	r.AddPattern("LOGLEVEL", `([Aa]lert|ALERT|[Tt]race|TRACE|[Dd]ebug|DEBUG|[Nn]otice|NOTICE|[Ii]nfo|INFO|[Ww]arn?(?:ing)?|WARN?(?:ING)?|[Ee]rr?(?:or)?|ERR?(?:OR)?|[Cc]rit?(?:ical)?|CRIT?(?:ICAL)?|[Ff]atal|FATAL|[Ss]evere|SEVERE|EMERG(?:ENCY)?|[Ee]merg(?:ency)?)`)
	r.AddPattern("MS_DELIM", `%{MS_DATE_TIME}\s%{MS_TZ}\s%{LOGLEVEL}`)
	return &LogReaderCompute{Regnet: r}
	//return &LogReaderCompute{}
}

func (reader *LogReaderCompute) Execute(arg Args) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		filePath := r.FormValue("path")
		//start
		rgexp, _ := regexp.Compile(`((\d\d[\/]*){3}[\s]*(\d\d[:]*){3}.(\d*)) ((?:[PMCEI][SD]T|UTC)|GMT-\d\d:\d\d) ([Aa]lert|ALERT|[Tt]race|TRACE|[Dd]ebug|DEBUG|[Nn]otice|NOTICE|[Ii]nfo|INFO|[Ww]arn?(?:ing)?|WARN?(?:ING)?|[Ee]rr?(?:or)?|ERR?(?:OR)?|[Cc]rit?(?:ical)?|CRIT?(?:ICAL)?|[Ff]atal|FATAL|[Ss]evere|SEVERE|EMERG(?:ENCY)?|[Ee]merg(?:ency)?)`)
		//end
		go reader.process(filePath, arg, rgexp)
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func (reader *LogReaderCompute) process(file string, arg Args, regexp *regexp.Regexp) {
	fmt.Println("parsing started at ::", time.Now())
	//initialize http handler and listen for POST message, get file name from request

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from a panic !", r)
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
			break
		}

		if regexp.Match(line) == true {
			if lg != nil && len(lg.Store) > 0 {
				//push it further
				packet := NewPacket()
				packet["log"] = lg
				packet["source"] = file
				/*logEntry := packet["log"].(*Logs)
				last := strings.TrimSpace(string(logEntry.Store[ (len(logEntry.Store) - 25) : len(logEntry.Store)]))
		
				fmt.Println("lastline<",last,">")
				if strings.EqualFold(last,"ThreadPoolImpl.java:127)") == false {
					fmt.Println("Jumbled:", string(logEntry.Store))
					fmt.Println("===========================================================")
					//fmt.Println("The String version:", logEntry.StoreInString)
					os.Exit(1)
				}*/
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