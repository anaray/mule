package mule

import (
	"github.com/anaray/compute"
)

type LogTypeIndentifierCompute struct{}

func LogTypeIndentifier() *LogTypeIndentifierCompute {
	return new(LogTypeIndentifierCompute)
}

func (identifier *LogTypeIndentifierCompute) String() string { return "compute.logtype_identifier" }

func (identifier *LogTypeIndentifierCompute) Execute(arg *compute.Args) {
	for {
		//packet := <-arg.Incoming
		//log := packet["log"].(*Logs)
		//log.Store

		//packetChannel <- packet
	}
}
