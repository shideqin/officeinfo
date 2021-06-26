package main

import (
	"flag"
	"fmt"
	"officeinfo/info"
	"os"
)

// VERSION 版本号
var VERSION = "0.0.1"

// HELP 帮助信息
var HELP = `
    version: ` + VERSION + `
    officeinfo

    image http[s]://domain/file[.jpg|.jpeg|.png|.gif]
    pdf  http[s]://domain/file.pdf [pageNum]
`

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Printf("%s\n", HELP)
		os.Exit(0)
	}

	//初始化office info
	oi := info.New()

	switch args[0] {
	case "image":
		resp, err := oi.GetImageInfo(args)
		if err != nil {
			fmt.Printf("image error: %s\n", err.Error())
		}
		fmt.Printf("%s", resp)
	case "pdf":
		resp, err := oi.GetPdfInfo(args)
		if err != nil {
			fmt.Printf("pdf error: %s\n", err.Error())
		}
		fmt.Printf("%s", resp)
	case "help":
		fmt.Printf("%s\n", HELP)
	default:
		fmt.Println("unsupported command : " + args[0])
		fmt.Println("use --help for more information")
	}
}
