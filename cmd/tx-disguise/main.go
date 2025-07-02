package main

import (
	"flag"
	"fmt"

	"tx-disguise/internal/futures"
	"tx-disguise/internal/tui"
)

const defaultFuturesCode = "TXF"
const version = "beta-0.1.0"

func usage() {
	fmt.Println(`
	Usage: tx-disguise [-v] [-h] [ -y | -z ]
		-v: show version  
		-h: show this help
	Symbol Options:
		-y: 小台 (MXF)  
		-z: 微台 (TMF)
	Example: 
		tx-disguise -y
	`)
}

func showVersion() {
	fmt.Printf("version: %s\n", version)
}

func main() {
	var (
		helpFlag    bool
		versionFlag bool
		yFlag       bool
		zFlag       bool
	)
	flag.BoolVar(&helpFlag, "h", false, "show help")
	flag.BoolVar(&versionFlag, "v", false, "show version")
	flag.BoolVar(&yFlag, "y", false, "小台 (MXF)")
	flag.BoolVar(&zFlag, "z", false, "微台 (TMF)")
	flag.Parse()

	if helpFlag {
		usage()
		return
	}
	if versionFlag {
		showVersion()
		return
	}
	futuresService := futures.NewService(defaultFuturesCode)
	if yFlag {
		futuresService.FuturesCode = "MXF"
	}
	if zFlag {
		futuresService.FuturesCode = "TMF"
	}
	view := tui.NewProgram(futuresService)
	if _, err := view.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		return
	}
}
