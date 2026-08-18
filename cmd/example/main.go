package main

import "tx-disguise/internal/futures"

func main() {
	service := futures.NewService("TXF")
	futuresPrice, _ := service.GetCurrentFuturesPrice()
	actualPrice, _ := service.GetCurrentActualPrice()
	println("Current Futures Price:", futuresPrice)
	println("Current Actual Price:", actualPrice)
}
