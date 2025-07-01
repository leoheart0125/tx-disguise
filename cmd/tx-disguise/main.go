package main

import (
	"flag"
	"fmt"
	"time"

	future "tx-disguise/internal/future"
)

var (
	isForcedClearScreen bool
	futuresCode         = "TXF"
	actualsCode         = "TXF-S"
	requestInterval     = 2 * time.Second
	version             = "0.7"
)

func usage() {
	fmt.Println(`
	Usage: tx-disguise [-r] [-v] [-h] [ -y | -z ]
		-r: clear price lines on every request
		-v: show version  
		-h: show this help
	Symbol Options:
		-y: 小台 (MXF)  
		-z: 微台 (TMF)
	Example: 
		tx-disguise -r  
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
		rFlag       bool
	)
	flag.BoolVar(&helpFlag, "h", false, "show help")
	flag.BoolVar(&versionFlag, "v", false, "show version")
	flag.BoolVar(&rFlag, "r", false, "clear price lines on every request")
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
	futureService := future.NewService(futuresCode)
	if yFlag {
		futureService.FuturesCode = "MXF"
	}
	if zFlag {
		futureService.FuturesCode = "TMF"
	}
	isForcedClearScreen = rFlag

	showVersion()
	fmt.Println()
	fmt.Println()
	fmt.Printf("%s %-11s %-21s | %-21s %s\n\n", "date", "", "Futures", "Actuals", "trash")

	for {
		fmt.Printf("[%s] %-21s | %-21s %s\n",
			time.Now().Format("01/02 15:04:05"),
			futureService.GetCurrentFuturePrice(),
			futureService.GetCurrentActualPrice(),
			"test",
		)
		time.Sleep(futureService.RequestInterval)
	}
}
