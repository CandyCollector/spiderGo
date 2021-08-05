package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
)

// 43：当前股价 57：股票编号 58：股票名字 86: 时间戳
// 46：今开 44：最高 51：涨停 47：成交量
// 60：昨收 45：最低 52：跌停 50：量比 48：成交额
type data struct {
	F43 float64 `json:"f43"`
	F44 float64 `json:"f44"`
	F45 float64 `json:"f45"`
	F47 int     `json:"f47"`
	F48 float64 `json:"f48"`
	F50 float64 `json:"f50"`
	// F51 string  `json:"f51"`
	// F52 float64 `json:"f52"`
	F57 string  `json:"f57"`
	F58 string  `json:"f58"`
	F60 float64 `json:"f60"`
	F86 int     `json:"f86"`
}

func main() {
	//连接 mongodb
	// dialInfo := &mgo.DialInfo{
	// 	Addrs:     []string{"120.77.*.*:27017"}, //远程(或本地)服务器地址及端口号
	// 	Direct:    false,
	// 	Timeout:   time.Second * 1,
	// 	Database:  "admin", //数据库
	// 	Source:    "admin",
	// 	Username:  "root",
	// 	Password:  "root",
	// 	PoolLimit: 4096, // Session.SetPoolLimit
	// }
	// session, err := mgo.DialWithInfo(dialInfo)
	// if err != nil {
	// 	panic(err)
	// }
	// defer session.Close()

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
		fmt.Println(string(r.Body))
		rel := string(r.Body)[54 : len(r.Body)-1]
		fmt.Println(rel)

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
	file, err := os.Open("number.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// This is our buffer now
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// fmt.Println("read lines:")
	// for _, line := range lines {
	// 	fmt.Println(line)
	// }
	for _, line := range lines {
		url := "http://push2.eastmoney.com/api/qt/stock/get?ut=fa5fd1943c7b386f172d6893dbfba10b&invt=2&fltt=2&fields=f43,f44,f51,f47,f60,f45,f52,f50,f48,f86,f58,f57&secid=1." + line + "&_=" + strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		fmt.Println(url)
		c.Visit(url)

	}

	// http://push2.eastmoney.com/api/qt/stock/get?ut=fa5fd1943c7b386f172d6893dbfba10b&invt=2&fltt=2&fields=f43,f57,f58,f169,f170,f46,f44,f51,f168,f47,f164,f163,f116,f60,f45,f52,f50,f48,f167,f117,f71,f161,f49,f530,f135,f136,f137,f138,f139,f141,f142,f144,f145,f147,f148,f140,f143,f146,f149,f55,f62,f162,f92,f173,f104,f105,f84,f85,f183,f184,f185,f186,f187,f188,f189,f190,f191,f192,f107,f111,f86,f177,f78,f110,f262,f263,f264,f267,f268,f250,f251,f252,f253,f254,f255,f256,f257,f258,f266,f269,f270,f271,f273,f274,f275,f127,f199,f128,f193,f196,f194,f195,f197,f80,f280,f281,f282,f284,f285,f286,f287,f292&secid=1.600519&_=1628157973294
	// url := "http://push2.eastmoney.com/api/qt/stock/get?ut=fa5fd1943c7b386f172d6893dbfba10b&invt=2&fltt=2&fields=f43,f44,f51,f47,f60,f45,f52,f50,f48,f86,f58,f57&secid=1.600519&_=1628144331000"
	// fmt.Println(url)
	// c.Visit(url)
	c.Wait()
	fmt.Printf("花费时间:%s", time.Since(t))
}
