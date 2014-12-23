A library to read log files, process them to extract/convert them to information and push to ElasticSearch for further analytics, reporting.

###### Getting Started
```
cd /myhome/mymule/
git clone https://github.com/anaray/mule.git

export GOPATH=/myhome/mymule/mule
export PATH=$GOPATH/bin:$PATH

# get dependencies - compute,regent,elastigo
go get github.com/anaray/compute/
go get github.com/anaray/regnet/ 
go get github.com/mattbaird/elastigo/

go install compute
go install mule
go install main 
    OR
go run src/main/mule.go 
mule #starts the mule process

```

Example : chaining a LogReader and ElasticSearchPusher
```
package main

import (
	"github.com/anaray/compute"
	"mule"
)

func main() {
	compute.Run(mule.LogReader(), mule.ElasticSearch())
}

```

Filters can added in between Source(compute.LogReader) & Sink(compute.ElasticSearch), to clean,filter & capture data
