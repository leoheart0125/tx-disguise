package futures

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type IService interface {
	GetCurrentFuturesPrice() (string, float64)
	GetCurrentActualPrice() (string, float64)
}

type Service struct {
	FuturesCode     string
	ActualsCode     string
	RequestInterval time.Duration
	httpClient      *http.Client
}

var _ IService = (*Service)(nil)

func NewService(futuresCode string) *Service {
	return &Service{
		FuturesCode:     futuresCode,
		ActualsCode:     "TXF-S",
		RequestInterval: 2 * time.Second,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *Service) apiRequest(symbolID string) (*Quote, error) {
	url := "https://mis.taifex.com.tw/futures/api/getChartData1M"
	payload := fmt.Sprintf(`{"SymbolID": "%s"}`, symbolID)
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Host", "mis.taifex.com.tw")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	var data RtData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data.RtData, nil
}

func (s *Service) GetCurrentFuturesPrice() (string, float64) {
	session := MarketSessionNow()
	var symbolID string
	contractCode := FuturesCurrentContractCode()
	symbolID = s.FuturesCode + contractCode
	switch session {
	case "regular":
		symbolID = symbolID + "-F"
	case "electronic":
		symbolID = symbolID + "-M"
	default:
		symbolID = ""
	}
	if symbolID == "" {
		return "-", 0
	}
	q, err := s.apiRequest(symbolID)
	if err != nil {
		return "-", 0
	}
	return ParseQuote(q), ParseFloat(q.Quote.CLastPrice)
}

func (s *Service) GetCurrentActualPrice() (string, float64) {
	q, err := s.apiRequest(s.ActualsCode)
	if err != nil {
		return "-", 0
	}
	return ParseQuote(q), ParseFloat(q.Quote.CLastPrice)
}
