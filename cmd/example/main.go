package main

import "tx-disguise/internal/futures"

func main() {
	service := futures.NewService("TXF")
	futuresPrice := service.GetCurrentFuturesPrice()
	actualPrice := service.GetCurrentActualPrice()
	println("Current Futures Price:", futuresPrice)
	println("Current Actual Price:", actualPrice)
}
