package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
	"gopkg.in/mgo.v2"
)

// 43：当前股价
// 46：今开 44：最高 51：涨停 47：成交量
// 60：昨收 45：最低 52：跌停 50：量比 48：成交额
// 86: 时间戳
type data struct {
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
}

func main() {
	//连接 mongodb
	dialInfo := &mgo.DialInfo{
		Addrs:     []string{"120.77.*.*:27017"}, //远程(或本地)服务器地址及端口号
		Direct:    false,
		Timeout:   time.Second * 1,
		Database:  "admin", //数据库
		Source:    "admin",
		Username:  "root",
		Password:  "root",
		PoolLimit: 4096, // Session.SetPoolLimit
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}
	defer session.Close()

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
		//Delay:     5 * time.Second,
	})

	// 处理接口返回 json 数据
	c.OnResponse(func(r *colly.Response) {
		//截取使用的数据
		rel := string(r.Body)[95 : len(r.Body)-3]
		// fmt.Println(rel)

		//解析 json 数据
		jue := new(data)
		err := json.Unmarshal([]byte(rel), jue)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(jue.F43)
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
