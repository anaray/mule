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
			//close(arg.Outgoing)

			//arg.Outgoing<-nil

			packet := NewPacket()
			packet["log"] = lg
			packet["source"] = file
			arg.Outgoing <- packet

			//emptyPacket := NewPacket()
			//emptyPacket["log"] = nil
			//emptyPacket["source"] = nil
			//emptyPacket:=nil
			//arg.Outgoing<-emptyPacket.(Packet)
			break
		}

		//_, err = reader.Regnet.MatchInText(string(line), "%{MS_DELIM}")
		/*res := regexp.FindAllString(string(line), -1)
		if res != nil {
			arg.Outgoing <- line
		}*/
		/*if err != nil {
			arg.Outgoing <- line
		}*/
		//fmt.Println("================================")
		//fmt.Println("LINE :",string(line))

		//all := regexp.FindAll(line, -1)
		/*for match := range all {
			fmt.Println("Match :",string(all[match]))
		}*/
		/*if all != nil {
			arg.Outgoing <- line
		}*/

		//message := string(line)
		if regexp.Match(line) == true {
			//count=count+1

			if len(lg.Store) > 0 {
				packet := NewPacket()
				packet["log"] = lg
				packet["source"] = file
				//fmt.Println("checking packet .....",string(lg.Store[0:8]) )
				if  string(lg.Store[0:8])!="07/09/14" {
					fmt.Println("FOUND ::", string(lg.Store[0:8]))
				}
				arg.Outgoing <- packet
			}
			lg = Logs{Store: line}

			//arg.Outgoing<-line
			//fmt.Println(" : MATCHED ")

			//}else{
			//fmt.Println(" : NOT MATCHED")
		} else {
			lg.Store = append(lg.Store[:], line[:]...)
		}
		//fmt.Println("================================")

		//if res != nil {
		//	arg.Outgoing <- line
		//}
	}
	//fmt.Println("Count at Reader::::::::::::::::::::::::::::::::::::::::::::::",count)

}
