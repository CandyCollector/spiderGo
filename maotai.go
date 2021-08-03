package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
)

// 43：当前股价
// 46：今开 44：最高 51：涨停 47：成交量
// 60：昨收 45：最低 52：跌停 50：量比 48：成交额
// 86: 时间戳
type juery struct {
	Rc   int `json:"rc"`
	Rt   int `json:"rt"`
	Svr  int `json:"svr"`
	Lt   int `json:"lt"`
	Full int `json:"full"`
	Data struct {
		F43 float64 `json:"f43"`
		F44 float64 `json:"f44"`
		F45 float64 `json:"f45"`
		F47 int     `json:"f47"`
		F48 float64 `json:"f48"`
		F50 float64 `json:"f50"`
		F51 float64 `json:"f51"`
		F52 float64 `json:"f52"`
		F60 float64 `json:"f60"`
		F86 int     `json:"f86"`
	} `json:"data"`
}

func main() {
	t := time.Now()
	c := colly.NewCollector(
		// 设置异步请求
		colly.Async(),
		// 开启 dubugger
		colly.Debugger(&debug.LogDebugger{}),
		// 域名过滤 支持正则
		colly.URLFilters(
			regexp.MustCompile("^(http://push2\\.eastmoney\\.com)/api"),
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
		jue := new(juery)
		err := json.Unmarshal(string(r.Body), jue)
    		if err != nil {
        		fmt.Println(err)
    		}
    		fmt.Println(jue)
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

	c.Visit("http://push2.eastmoney.com/api/qt/stock/get?ut=fa5fd1943c7b386f172d6893dbfba10b&invt=2&fltt=2&fields=f43,f44,f51,f47,f60,f45,f52,f50,f48,f86&secid=1.600519&cb=jQuery11240521556968571107_1627939429631&_=1627939429632")
	c.Wait()
	fmt.Printf("花费时间:%s", time.Since(t))
}
