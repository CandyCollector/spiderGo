package main

import (
	"fmt"
	"os"
)

func main() {
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
	var num string = string(buffer[:])
	for i := 0; i < len(num); i++ {
		fmt.Printf(num)
	}
}
