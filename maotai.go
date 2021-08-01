package main

import (
	"fmt"
	"time"

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
		// colly.URLFilters(
		// 	regexp.MustCompile("^(http://quote\\.eastmoney\\.com)/sh\\d{1,6}\\.html"),
		// ),
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
	// #price9
	// //*[@id="day"]  #price9

	// c.OnHTML("body", func(e *colly.HTMLElement) {

	// 	// test := e.DOM.Find("/html/body/div[10]/div[2]/span[1]").Text()
	// 	// fmt.Println("test:", test)
	// 	e.DOM.Each(func(i int, selection *goquery.Selection) {
	// 		// 	// name := selection.Find("#name").Text()
	// 		// 	// code := selection.Find("#code").Text()
	// 		// 	// time := selection.Find("#day").Text()
	// 		// 	// price := selection.Find("#gt1").Text()
	// 		fmt.Println(selection.Text())
	// 	})

	// })

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(string(r.Body))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visting", r.URL)

	})

	c.OnError(func(response *colly.Response, err error) {
		fmt.Println(err)
	})

	// 遍历股票 URL 地址
	// file, err := os.Open("number.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer file.Close()

	// scanner := bufio.NewScanner(file)
	// scanner.Split(bufio.ScanLines)

	// // This is our buffer now
	// var lines []string

	// for scanner.Scan() {
	// 	lines = append(lines, scanner.Text())
	// }

	// // fmt.Println("read lines:")
	// // for _, line := range lines {
	// // 	fmt.Println(line)
	// // }
	// for _, line := range lines {
	// 	url := "http://quote.eastmoney.com/sh" + line + ".html"
	// 	// fmt.Println(url)
	// 	c.Visit(url)

	// }

	// c.Visit("http://quote.eastmoney.com/sh600519.html")

	c.Visit("http://push2.eastmoney.com/api/qt/stock/get?ut=fa5fd1943c7b386f172d6893dbfba10b&invt=2&fltt=2&fields=f43,f57,f58,f169,f170,f46,f44,f51,f168,f47,f164,f163,f116,f60,f45,f52,f50,f48,f167,f117,f71,f161,f49,f530,f135,f136,f137,f138,f139,f141,f142,f144,f145,f147,f148,f140,f143,f146,f149,f55,f62,f162,f92,f173,f104,f105,f84,f85,f183,f184,f185,f186,f187,f188,f189,f190,f191,f192,f107,f111,f86,f177,f78,f110,f262,f263,f264,f267,f268,f250,f251,f252,f253,f254,f255,f256,f257,f258,f266,f269,f270,f271,f273,f274,f275,f127,f199,f128,f193,f196,f194,f195,f197,f80,f280,f281,f282,f284,f285,f286,f287,f292&secid=1.600519&cb=jQuery112406636451664591729_1627800779195&_=1627800779211")
	c.Wait()
	fmt.Printf("花费时间:%s", time.Since(t))
}
