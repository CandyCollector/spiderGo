package main

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
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
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	//获取页面数据
	// body > div.qphox.header_title.mb7  <h2 class="header-title-h2 fl" id="name">贵州茅台</h2>
	c.OnHTML("#div h2", func(e *colly.HTMLElement) {
		e.DOM.Each(func(i int, selection *goquery.Selection) {
			name := selection.Find("#name").Text()
			fmt.Print("%s", name)
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
	fmt.Println("输出测试")
	fmt.Printf("花费时间:%s", time.Since(t))
}
