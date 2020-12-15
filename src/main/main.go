package main

import (
	"fmt"
	"os"
	"rssreader/src/router"
)

func main() {
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "start":
			if len(os.Args) == 2 {
				router.Run("config.json")
			}
			if len(os.Args) == 3 {
				router.Run(os.Args[2])
			}
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
