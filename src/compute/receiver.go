package compute

import (
	"fmt"
	"time"
)

type StdoutPlugin struct{}

func Stdout() *StdoutPlugin {
	return new(StdoutPlugin)
}

func (p *StdoutPlugin) Execute(arg Args) {
	//var msg int = 0
	//var log *Logs
	for {
		//count = count + 1
		packet := <-arg.Incoming
		if packet == nil {
			break
		}
		//fmt.Println( ">> ", packet)
		//msg=msg+1
		//fmt.Println("Count at Receiver :", packet)
	}

	/*select {
	case packet := <-arg.Incoming:

	case <-arg.Done

	}*/

	fmt.Println("End at :", time.Now())
}
