package future

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	FuturesCode     = "TXF"
	ActualsCode     = "TXF-S"
	RequestInterval = 2 * time.Second
)

func MarketSessionNow() string {
	nowMinutes := getNowTotalMinutes()
	if nowMinutes >= 525 && nowMinutes <= 825 {
		return "regular"
	}
	if nowMinutes <= 300 || nowMinutes >= 900 {
		return "electronic"
	}
	return "closed"
}

func getNowTotalMinutes() int {
	now := time.Now()
	return now.Hour()*60 + now.Minute()
}

func SelfAPIRequest(symbolID string) (*Quote, error) {
	url := "https://mis.taifex.com.tw/futures/api/getChartData1M"
	payload := fmt.Sprintf(`{"SymbolID": "%s"}`, symbolID)
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Host", "mis.taifex.com.tw")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data RtData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data.RtData, nil
}

func FuturesIsThisMonthSettled() bool {
	now := time.Now()
	year, month, day := now.Date()
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	weekdayOf1st := int(firstOfMonth.Weekday())
	dateOfFirstWednesday := (11 - weekdayOf1st) % 7
	settleDay := dateOfFirstWednesday + 14
	return day > settleDay
}

func FuturesCurrentContractCode() string {
	now := time.Now()
	year, month, _ := now.Date()
	deltaMonth := 0
	if month == 12 {
		monthHex := 0x76
		if FuturesIsThisMonthSettled() {
			year++
			monthHex = 0x65
		}
		return fmt.Sprintf("%c%d", monthHex, year%10)
	}
	if FuturesIsThisMonthSettled() {
		deltaMonth = 1
	}
	monthHex := 64 + int(month) + deltaMonth
	return fmt.Sprintf("%c%d", monthHex, year%10)
}

func FutureGetCurrentQuote() *Quote {
	session := MarketSessionNow()
	var symbolID string
	switch session {
	case "regular":
		symbolID = FuturesCode + FuturesCurrentContractCode() + "-F"
	case "electronic":
		symbolID = FuturesCode + FuturesCurrentContractCode() + "-M"
	default:
		symbolID = ""
	}
	if symbolID == "" {
		return nil
	}
	q, err := SelfAPIRequest(symbolID)
	if err != nil {
		return nil
	}
	return q
}

func ActualsGetCurrentQuote() *Quote {
	q, err := SelfAPIRequest(ActualsCode)
	if err != nil {
		return nil
	}
	return q
}

func SelfGetPrice(q *Quote) string {
	if q == nil {
		return "-"
	}
	lastPrice := ParseInt(q.Quote.CLastPrice)
	refPrice := ParseInt(q.Quote.CRefPrice)
	highPrice := ParseInt(q.Quote.CHighPrice)
	lowPrice := ParseInt(q.Quote.CLowPrice)
	priceDiff := lastPrice - refPrice
	diffStr := fmt.Sprintf("%+d", priceDiff)
	return fmt.Sprintf("%d %s (%d, %d)", lastPrice, diffStr, lastPrice-lowPrice, highPrice-lastPrice)
}

func ParseInt(s string) int {
	if i := strings.Index(s, "."); i >= 0 {
		s = s[:i]
	}
	var v int
	if _, err := fmt.Sscanf(s, "%d", &v); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing integer from string '%s': %v\n", s, err)
		return 0
	}
	return v
}

func FakeInfo() string {
	cmd := exec.Command("top", "-l", "1")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "CPU") {
			return "==> " + line
		}
	}
	return ""
}

func ClearScreen(isForcedClearScreen bool) {
	if isForcedClearScreen {
		fmt.Print("\033[1A\033[K")
	}
}
