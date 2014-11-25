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
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	// make a read buffer
	r := bufio.NewReader(fi)
	//var lg Logs
	var count int
	for {
		// read a chunk
		_, err := r.ReadString('\n')
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
