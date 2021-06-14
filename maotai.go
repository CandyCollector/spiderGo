package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
	
)

func main(){
	t := time.Now()
	c := colly.NewCollector(
		// 设置异步请求
		colly.Async(),
		// 开启 dubugger
		colly.Debugger(&debug.LogDebugger{}),
		// 域名过滤 支持正则
		// https://finance.sina.com.cn/realstock/company/sh600519/nc.shtml
		colly.URLFilters(
			regexp.MustCompile("https://finance\\.sina\\.com\\.cn/realstock/company/^sh\\d{1,6}/\\.nc\\.shtml"),	
	))
	//使用扩展插件 
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	c.OnHTML("a[href]",func(e *colly.HTMLElement){
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request){
		fmt.Println("Visting",r.URL)
	})

	c.OnError(func(response *colly.Response, err error) {
        fmt.Println(err)
    })


	c.Visit("https://finance.sina.com.cn/realstock/company")
	fmt.Printf("花费时间:%s",time.Since(t))
}
