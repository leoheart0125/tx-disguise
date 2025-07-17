package futures

import (
	"testing"
	"time"
)

func TestParseInt(t *testing.T) {
	cases := []struct {
		input    string
		expected int
	}{
		{"123", 123},
		{"", 0},
		{"abc", 0},
		{"42.5", 42},
		{"-7", -7},
	}
	for _, c := range cases {
		if got := ParseInt(c.input); got != c.expected {
			t.Errorf("ParseInt(%q) = %d; want %d", c.input, got, c.expected)
		}
	}
}

func TestParseQuote(t *testing.T) {
	q := &Quote{}
	q.Quote.CLastPrice = "100"
	q.Quote.CRefPrice = "90"
	q.Quote.CHighPrice = "110"
	q.Quote.CLowPrice = "80"
	result := ParseQuote(q)
	if result == "-" || result == "" {
		t.Errorf("ParseQuote failed, got %q", result)
	}
	if ParseQuote(nil) != "-" {
		t.Errorf("ParseQuote(nil) should return '-' ")
	}
}

func TestMarketSessionNow(t *testing.T) {
	// This test is time-dependent, so just check output is one of the expected values
	valid := map[string]bool{"regular": true, "electronic": true, "closed": true}
	if !valid[MarketSessionNow()] {
		t.Errorf("MarketSessionNow returned unexpected value: %q", MarketSessionNow())
	}
}

func TestFuturesCurrentContractCode(t *testing.T) {
	code := FuturesCurrentContractCode()
	if len(code) != 2 {
		t.Errorf("FuturesCurrentContractCode returned %q, want 2 chars", code)
	}
}

func TestFuturesIsThisMonthSettled(t *testing.T) {
	// This test is time-dependent, so just check it returns a bool
	_ = FuturesIsThisMonthSettled()
}

func TestMarketSessionAt(t *testing.T) {
	cases := []struct {
		hour, min int
		expected  string
	}{
		{8, 45, "regular"},     // 8:45 AM
		{5, 0, "electronic"},   // 5:00 AM
		{15, 30, "electronic"}, // 3:30 PM
		{13, 50, "closed"},     // 1:50 PM
	}
	for _, c := range cases {
		tm := time.Date(2025, 7, 17, c.hour, c.min, 0, 0, time.Local)
		if got := MarketSessionAt(tm); got != c.expected {
			t.Errorf("MarketSessionAt(%02d:%02d) = %q; want %q", c.hour, c.min, got, c.expected)
		}
	}
}

func TestFuturesIsThisMonthSettledAt(t *testing.T) {
	// Test a date before settlement
	beforeSettle := time.Date(2025, 7, 16, 8, 0, 0, 0, time.Local)
	if FuturesIsThisMonthSettledAt(beforeSettle) {
		t.Error("Expected not settled before settlement time")
	}
	// Test a date after settlement
	afterSettle := time.Date(2025, 7, 16, 14, 0, 0, 0, time.Local)
	if !FuturesIsThisMonthSettledAt(afterSettle) {
		t.Error("Expected settled after settlement time")
	}
}

func TestFuturesCurrentContractCodeAt(t *testing.T) {
	// Test December rollover
	dec := time.Date(2025, 12, 20, 10, 0, 0, 0, time.Local)
	code := FuturesCurrentContractCodeAt(dec)
	if len(code) != 2 {
		t.Errorf("FuturesCurrentContractCodeAt returned %q, want 2 chars", code)
	}
	// Test normal month
	jul := time.Date(2025, 7, 10, 10, 0, 0, 0, time.Local)
	code2 := FuturesCurrentContractCodeAt(jul)
	if len(code2) != 2 {
		t.Errorf("FuturesCurrentContractCodeAt returned %q, want 2 chars", code2)
	}
}
