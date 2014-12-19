A library to read log files, process them to extract/convert them to information and push to ElasticSearch for further analytics, reporting.

###### Getting Started
```
cd /myhome/mymule/
git clone https://github.com/anaray/mule.git
export GOPATH=/myhome/mymule/mule
export PATH=$GOPATH/bin:$PATH
go get github.com/anaray/regnet/ 
go get github.com/mattbaird/elastigo/

go install compute
go install mule
mule //starts the mule process

```

Example : chaining a LogReader and ElasticSearchPusher
```
package main

import ("compute")

func main() {
	compute.Run(compute.LogReader(), compute.ElasticSearch()))
}
```

Filters can added in between Source(compute.LogReader) & Sink(compute.ElasticSearch), to clean,filter & capture data
