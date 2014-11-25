package compute

import (
	"bufio"
	"io"
	"os"
	//"io/ioutil"
	"fmt"
	"github.com/anaray/regnet"
	"net/http"
	"regexp"
	"time"
)

//https://groups.google.com/forum/#!topic/golang-nuts/sAwDldpkMGQ
type LogReaderCompute struct {
	Regnet *regnet.Regnet
}

type Logs struct {
	Store []byte
}

func LogReader() *LogReaderCompute {
	r, _ := regnet.New()
	r.AddPattern("MS_DATE_TIME", `((\d\d[\/]*){3}[\s]*(\d\d[:]*){3}.\d\d\d)`)
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
		rgexp, _ := regexp.Compile(`((\d\d[\/]*){3}[\s]*(\d\d[:]*){3}.\d\d\d) ((?:[PMCEI][SD]T|UTC)|GMT-\d\d:\d\d) ([Aa]lert|ALERT|[Tt]race|TRACE|[Dd]ebug|DEBUG|[Nn]otice|NOTICE|[Ii]nfo|INFO|[Ww]arn?(?:ing)?|WARN?(?:ING)?|[Ee]rr?(?:or)?|ERR?(?:OR)?|[Cc]rit?(?:ical)?|CRIT?(?:ICAL)?|[Ff]atal|FATAL|[Ss]evere|SEVERE|EMERG(?:ENCY)?|[Ee]merg(?:ency)?)`)
		//end
		go reader.process(filePath, arg, rgexp)
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func (reader *LogReaderCompute) process(file string, arg Args, regexp *regexp.Regexp) {
	fmt.Println("Starting parsing at ::", time.Now())
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
	var lg Logs
	//var count int = 0
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
			if len(lg.Store) > 0 {
				packet := NewPacket()
				packet["log"] = lg
				packet["source"] = file
				arg.Outgoing <- packet
			}
			lg = Logs{Store: line}
		} else {
			lg.Store = append(lg.Store[:], line[:]...)
		}
	}
}
