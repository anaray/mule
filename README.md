A simple example: chaining LogReader -> ElasticSearchPusher
```
package main

import ("compute")

func main() {
	compute.Run(compute.LogReader(), compute.ElasticSearch()))
}
```

Filters can added in between Source(compute.LogReader) & Sink(compute.ElasticSearch), to clean,filter & capture data
