package futures

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseInt(s string) int {
	// avoid parsing errors from empty strings or non-numeric values
	if s == "" {
		return 0
	}
	if i := strings.Index(s, "."); i >= 0 {
		s = s[:i]
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

func FuturesIsThisMonthSettledAt(now time.Time) bool {
	year, month, day := now.Date()
	nowMinutes := now.Hour()*60 + now.Minute()
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	weekdayOf1st := int(firstOfMonth.Weekday())
	dateOfFirstWednesday := (11 - weekdayOf1st) % 7
	settleDay := dateOfFirstWednesday + 14
	if day == settleDay && (nowMinutes <= 825) {
		return false
	}
	return day > settleDay || (day == settleDay && nowMinutes > 825)
}

func FuturesIsThisMonthSettled() bool {
	return FuturesIsThisMonthSettledAt(time.Now())
}

func FuturesCurrentContractCodeAt(now time.Time) string {
	year, month, _ := now.Date()
	deltaMonth := 0
	isMonthSettled := FuturesIsThisMonthSettledAt(now)
	if month == 12 {
		if isMonthSettled {
			year++
			month = 1
		}
		return fmt.Sprintf("%c%d", month, year%10)
	}
	if isMonthSettled {
		deltaMonth = 1
	}
	monthHex := 64 + int(month) + deltaMonth
	return fmt.Sprintf("%c%d", monthHex, year%10)
}

func FuturesCurrentContractCode() string {
	return FuturesCurrentContractCodeAt(time.Now())
}

func ParseQuote(q *Quote) string {
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

func MarketSessionAt(now time.Time) string {
	nowMinutes := now.Hour()*60 + now.Minute()
	if nowMinutes >= 525 && nowMinutes <= 825 {
		return "regular"
	}
	if nowMinutes <= 300 || nowMinutes >= 900 {
		return "electronic"
	}
	return "closed"
}

func MarketSessionNow() string {
	return MarketSessionAt(time.Now())
}
