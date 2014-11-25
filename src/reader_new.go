package main

import (
	"bufio"
	"io"
	//"log"
	"os"
	"time"
	//"io/ioutil"
	"fmt"
)

func readLines(path string) {
	startTime := time.Now()
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	// make a read buffer
	reader := bufio.NewReaderSize(file, 1024*1024)
	//var lg Logs
	var count int
	for {
		// read a chunk
		_, err := reader.ReadSlice('\n')
		//fmt.Println(string(line))
		count=count+1
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			endTime := time.Now()
			fmt.Println("completed in ", endTime.Sub(startTime))
			break
		}

		//fmt.Println(line)

	}
}

func main() {
	readLines("/home/msi/Desktop/applog_covidien.log")
}
