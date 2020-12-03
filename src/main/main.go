package main

import (
	"fmt"
	"os"
	"rssreader/src/router"
)

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "start":
			router.Run()
			break
		case "test":
			fmt.Println("OK")
			return
		default:
			fmt.Println("参数错误")
		}
	}
	return
}
