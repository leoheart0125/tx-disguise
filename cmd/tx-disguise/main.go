package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	future "tx-disguise/internal/txdisguise"
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
	if yFlag {
		future.FuturesCode = "MXF"
	}
	if zFlag {
		future.FuturesCode = "TMF"
	}
	isForcedClearScreen = rFlag

	showVersion()
	fmt.Println()
	futQuote := future.FutureGetCurrentQuote()
	if futQuote != nil {
		fmt.Fprintf(os.Stderr, "%s (%s)\n", futQuote.DispCName, future.MarketSessionNow())
	}
	fmt.Println()
	fmt.Printf("%s %-11s %-21s | %-21s %s\n\n", "date", "", "Futures", "Actuals", "trash")

	for {
		future.ClearScreen(isForcedClearScreen)
		fmt.Printf("[%s] %-21s | %-21s %s\n",
			time.Now().Format("01/02 15:04:05"),
			future.SelfGetPrice(future.FutureGetCurrentQuote()),
			future.SelfGetPrice(future.ActualsGetCurrentQuote()),
			future.FakeInfo(),
		)
		time.Sleep(future.RequestInterval)
	}
}
