package main

import ("compute")

func main() {
	compute.Run(compute.LogReader(), compute.ElasticSearch()) //compute.Stdout())
}