package main

import ("compute")

func main() {
	//compute.SetConfiguration("/home/msi/Desktop/test_mule.toml")
	compute.Run(compute.LogReader(), compute.ElasticSearch()) //compute.Stdout())
}