A library to read log files, process them to extract/convert them to information and push to ElasticSearch for further analytics, reporting.

Exmaple : chaining a LogReader and ElasticSearchPusher
```
package main

import ("compute")

func main() {
	compute.Run(compute.LogReader(), compute.ElasticSearch()))
}
```

Filters can added in between Source(compute.LogReader) & Sink(compute.ElasticSearch), to clean,filter & capture data
