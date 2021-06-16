package main

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
)

func main() {
	t := time.Now()
	c := colly.NewCollector(
		// 设置异步请求
		colly.Async(),
		// 开启 dubugger
		colly.Debugger(&debug.LogDebugger{}),
		// 域名过滤 支持正则
		// https://finance.sina.com.cn/realstock/company/sh600519/nc.shtml
		// colly.URLFilters(
		// 	regexp.MustCompile("^(http://quote\\.eastmoney\\.com)/sh\\d{1,6}\\.html"),
		// )
	)
	//使用扩展插件
	// extensions.RandomUserAgent(c)
	// extensions.Referer(c)

	//获取页面数据
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.DOM.Each(func(i int, selection *goquery.Selection) {
			// #name
			name := selection.Find("header-title-h1 f1").First().Text()

			fmt.Println("输出测试")
			fmt.Println(name)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visting", r.URL)

	})

	c.OnError(func(response *colly.Response, err error) {
		fmt.Println(err)
	})

	// http://quote.eastmoney.com/sh600519.html

	c.Visit("http://quote.eastmoney.com/sh600519.html")
	fmt.Printf("花费时间:%s", time.Since(t))
}
