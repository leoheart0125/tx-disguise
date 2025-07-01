package main

import (
	"tx-disguise/internal/future"
	"tx-disguise/internal/tui"
)

func main() {
	p := tui.NewProgram(future.NewService("TXF"))
	p.Run()
}
