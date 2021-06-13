package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
)

func main() {
	url := "https://finance.sina.com.cn/realstock/company/sh600519/nc.shtml"

	// Instantiate default collector
	c := colly.NewCollector(
		// Turn on asynchronous requests
		colly.Async(),
		// Attach a debugger to the collector
		colly.Debugger(&debug.LogDebugger{}),
	)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*https.*",
		Parallelism: 2,
		//Delay:      5 * time.Second,
	})

	// Start scraping in five threads on 
	for i := 0; i < 5; i++ {
		c.Visit(fmt.Sprintf("%s?n=%d", url, i))
	}
	// Wait until threads are finished
	c.Wait()
}