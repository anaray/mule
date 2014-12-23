package main

import (
	"github.com/anaray/compute"
	"mule"
)

func main() {
	//compute.SetConfiguration("/home/msi/Desktop/test_mule.toml")
	compute.Run(mule.LogReader(), mule.ElasticSearch()) //compute.Stdout())
}
