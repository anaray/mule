package main

import (
	"bufio"
	"io"
	"os"
	"time"
	// "bytes"
	"fmt"
)

// Read a whole file into the memory and store it as array of lines
func readLines(path string) {
	startTime := time.Now()
	file, _ := os.Open(path)
	reader := bufio.NewReaderSize(file, 1024*1024)
	
	//reader := bufio.NewReader(file)
	//var seq []byte
	var count int
	for {
		_, err := reader.ReadSlice('\n')
		//seq = append(seq,line[0:len(line)-1]...)
		count=count+1
		if err != nil {
			if err == io.EOF {
				endTime := time.Now()
				fmt.Println("completed in ", endTime.Sub(startTime))
				//fmt.Println("completed ", count)
				break
			} else {
				panic(err)
			}
		}
	}
	
}

func main() {
	//readLines("/home/msi/logstash/log/anaray/trakr_111/applog_new.log")
	readLines("/home/msi/Desktop/applog_covidien.log")
	//fmt.Println(string(lines))
	/*for _, line := range lines {
	    fmt.Println(line)
	}*/
}
