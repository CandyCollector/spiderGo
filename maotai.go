package main

import (
	"fmt"

	"github.com/gocolly/colly/"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"regexp"
)

func main(){
	t := time.Now()
	number := 1
	c := colly.NewCollector(
		// 设置异步请求
		colly.Async(),
		// 开启 dubugger
		colly.Debugger(&debug.LogDebugger{}),
		// 域名过滤 支持正则
		// https://finance.sina.com.cn/realstock/company/sh600519/nc.shtml
		colly.URLFilters(
			regexp,MustCompile("https://finance\\.sina\\.com\\.cn/realstock/company/^sh\\d{1,6}/\\.nc\\.shtml"),	
	))

	extensions.RandomUserAgent(c)
	extensions.Refer(c)

	c.OnHTML("a[href]",func(e *colly.HTMLElement){
		e.request. Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request){
		fmt.Println("Visting",r.url)
	})

	c.Visit("https://finance.sina.com.cn/realstock/company")

}
