package future

type Quote struct {
	Quote struct {
		CLastPrice string `json:"CLastPrice"`
		CRefPrice  string `json:"CRefPrice"`
		CHighPrice string `json:"CHighPrice"`
		CLowPrice  string `json:"CLowPrice"`
	} `json:"Quote"`
	DispCName string `json:"DispCName"`
}

type RtData struct {
	RtData Quote `json:"RtData"`
}
