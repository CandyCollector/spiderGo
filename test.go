package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	t1 := time.Now().UnixNano() / 1e6
	fmt.Println(t1)
	// c := time.Unix(t1, 0)
	url := "http://push2.eastmoney.com/api/qt/stock/get?ut=" + strconv.FormatInt(t1, 10)
	fmt.Println(url)
}
