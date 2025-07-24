package main

import (
	"flag"
	"fmt"

	"tx-disguise/internal/fakeinfo"
	"tx-disguise/internal/futures"
	"tx-disguise/internal/tui"
)

const defaultFuturesCode = "TXF"
const version = "v0.1.0"

func usage() {
	fmt.Println(`
	Usage: txd [-v] [-h] [ -y | -z ]
		-v: show version  
		-h: show this help
	Symbol Options:
		-y: 小台 (MXF)  
		-z: 微台 (TMF)
	Fake Info Options:
		-f: fake info service type (top or pure)
			Default is "pure" which uses a pure fake system info service.
			"top" uses the system's top command to get process info only on Apple Silicon.
	Example:
		txd 
		txd -y -f top
	`)
}

func showVersion() {
	fmt.Printf("version: %s\n", version)
}

func main() {
	var (
		helpFlag     bool
		versionFlag  bool
		yFlag        bool
		zFlag        bool
		fakeInfoFlag string
	)
	flag.BoolVar(&helpFlag, "h", false, "show help")
	flag.BoolVar(&versionFlag, "v", false, "show version")
	flag.BoolVar(&yFlag, "y", false, "小台 (MXF)")
	flag.BoolVar(&zFlag, "z", false, "微台 (TMF)")
	flag.StringVar(&fakeInfoFlag, "f", "pure", "fake info service type (top or pure)")
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
	fakeInfoService := fakeinfo.FackInfoFactory(fakeInfoFlag)
	view := tui.NewProgram(futuresService, fakeInfoService)
	if _, err := view.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		return
	}
}
