package main

import (
	"fmt"
	"tx-disguise/internal/fakeinfo"
)

func main() {
	fakeinfoService := fakeinfo.NewPureFakeSystemInfoService()
	fakeInfo, err := fakeinfoService.GetFakeInfo()
	if err != nil {
		fmt.Println("Error getting fake info:", err)
		return
	}
	for _, info := range fakeInfo {
		fmt.Println(info)
	}
}
