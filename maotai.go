package main

import (
	"fmt"
	"os"
	"regexp"
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
		//colly.Async(),
		// 开启 dubugger
		colly.Debugger(&debug.LogDebugger{}),
		// 域名过滤 支持正则
		colly.URLFilters(
			regexp.MustCompile("^(http://quote\\.eastmoney\\.com)/sh\\d{1,6}\\.html"),
		),
	)
	//使用扩展插件
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5,
		//Delay:      5 * time.Second,
	})

	//获取页面数据
	// body > div.qphox.header_title.mb7  <h2 class="header-title-h2 fl" id="name">贵州茅台</h2>

	c.OnHTML("body > div.qphox.header_title.mb7 ", func(e *colly.HTMLElement) {
		e.DOM.Each(func(i int, selection *goquery.Selection) {
			name := selection.Find("h2").Text()
			if len(name) != 0 {
				fmt.Println(name, "：", e.Request.URL)
			}

		})

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visting", r.URL)

	})

	c.OnError(func(response *colly.Response, err error) {
		fmt.Println(err)
	})

	// 遍历股票 URL 地址
	file, err := os.Open("number.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("bytes read: ", bytesread)

	for i := 0; i <= 14103; i++ {
		var num = string(buffer)
		url := "http://quote.eastmoney.com/sh" + num + ".html"
		// fmt.Println(url)
		c.Visit(url)

	}

	// c.Visit("http://quote.eastmoney.com/sh000001.html")
	c.Wait()
	fmt.Printf("花费时间:%s", time.Since(t))
}
