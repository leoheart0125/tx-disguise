package tui

import (
	"fmt"
	"os"
	"strings"
	"time"
	"tx-disguise/internal/fakeinfo"
	"tx-disguise/internal/futures"
	"tx-disguise/internal/shared"

	"github.com/NimbleMarkets/ntcharts/canvas/runes"
	"github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

const totalTabs = 3

type model struct {
	futuresService  futures.IService
	fakeInfoService fakeinfo.IService
	fakeInfo        []string
	futures         string
	futuresHistory  *shared.RingBuffer[string]
	height          int
	width           int
	currentTab      int // 0: fakeInfo, 1: futuresHistory, 2: chart
	historyScroll   int // scroll position for futuresHistory

	// cached price values from the 2s ticker, reused by chart
	lastFuturesPrice float64
	lastActualsPrice float64

	// chart tab
	futuresChart       timeserieslinechart.Model
	actualsChart       timeserieslinechart.Model
	futuresChartInited bool
	actualsChartInited bool
	chartStartTime     time.Time     // when chart data started
	chartWindowDur     time.Duration // sliding window duration
}

type (
	fakeInfoMsg       []string
	futuresHistoryMsg *shared.RingBuffer[string]
	chartTickMsg      struct{}
)

type futuresMsg struct {
	display      string
	futuresPrice float64
	actualsPrice float64
}

func (m model) fakeInfoTicker() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		fakeInfo, err := m.fakeInfoService.GetFakeInfo()
		if err != nil {
			return fakeInfoMsg([]string{"[ERROR] " + err.Error()})
		}
		return fakeInfoMsg(fakeInfo)
	})
}

func (m model) futuresTicker() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		fStr, fVal := m.futuresService.GetCurrentFuturesPrice()
		aStr, aVal := m.futuresService.GetCurrentActualPrice()
		return futuresMsg{
			display: fmt.Sprintf("[%s] %-21s | %-21s ",
				time.Now().Format("01/02 15:04:05"), fStr, aStr),
			futuresPrice: fVal,
			actualsPrice: aVal,
		}
	})
}

func (m model) futuresHistoryTicker() tea.Cmd {
	return tea.Tick(1*time.Minute, func(t time.Time) tea.Msg {
		m.futuresHistory.Push(m.futures)
		return futuresHistoryMsg(m.futuresHistory)
	})
}

func (m model) chartTicker() tea.Cmd {
	return tea.Tick(1*time.Minute, func(t time.Time) tea.Msg {
		return chartTickMsg{}
	})
}

// chartFirstTickCmd fires quickly after startup to get initial chart data
func chartFirstTickCmd() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return chartTickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.fakeInfoTicker(), m.futuresTicker(), m.futuresHistoryTicker(), chartFirstTickCmd())
}

// initChart sets Y range and recalculates graph sizes for proper label widths
func initChart(chart *timeserieslinechart.Model, price float64) {
	chart.SetYRange(price-200, price+200)
	chart.SetViewYRange(price-200, price+200)
	chart.UpdateGraphSizes()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fakeInfoMsg:
		m.fakeInfo = msg
		return m, m.fakeInfoTicker()

	case futuresMsg:
		m.futures = msg.display
		m.lastFuturesPrice = msg.futuresPrice
		m.lastActualsPrice = msg.actualsPrice
		return m, m.futuresTicker()

	case futuresHistoryMsg:
		m.futuresHistory = msg
		return m, m.futuresHistoryTicker()

	case chartTickMsg:
		now := time.Now()
		futuresPrice := m.lastFuturesPrice
		actualsPrice := m.lastActualsPrice

		// Record start time on first data from either chart
		if m.chartStartTime.IsZero() && (futuresPrice > 0 || actualsPrice > 0) {
			m.chartStartTime = now
		}

		if futuresPrice > 0 {
			if !m.futuresChartInited {
				initChart(&m.futuresChart, futuresPrice)
				m.futuresChartInited = true
			}
			m.futuresChart.Push(timeserieslinechart.TimePoint{Time: now, Value: futuresPrice})
		}
		if actualsPrice > 0 {
			if !m.actualsChartInited {
				initChart(&m.actualsChart, actualsPrice)
				m.actualsChartInited = true
			}
			m.actualsChart.Push(timeserieslinechart.TimePoint{Time: now, Value: actualsPrice})
		}

		// Sliding window: only activate after window duration has elapsed
		if !m.chartStartTime.IsZero() && (m.futuresChartInited || m.actualsChartInited) {
			windowEnd := m.chartStartTime.Add(m.chartWindowDur)
			if now.After(windowEnd) {
				// Slide: show [now - windowDur, now]
				viewStart := now.Add(-m.chartWindowDur)
				if m.futuresChartInited {
					m.futuresChart.SetViewTimeRange(viewStart, now)
				}
				if m.actualsChartInited {
					m.actualsChart.SetViewTimeRange(viewStart, now)
				}
			}
			// Redraw
			if m.futuresChartInited {
				m.futuresChart.DrawBraille()
			}
			if m.actualsChartInited {
				m.actualsChart.DrawBraille()
			}
		}
		return m, m.chartTicker()

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "tab":
			m.currentTab = (m.currentTab + 1) % totalTabs
			return m, nil
		case "up":
			if m.currentTab == 1 {
				history := m.futuresHistory.GetAll()
				maxScroll := max(len(history)-max(m.height-3, 0), 0)
				if m.historyScroll < maxScroll {
					m.historyScroll++
				}
			}
			return m, nil
		case "down":
			if m.currentTab == 1 && m.historyScroll > 0 {
				m.historyScroll--
			}
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.currentTab {
	case 0:
		return m.fakeInfoView()
	case 1:
		return m.historyView()
	case 2:
		return m.chartView()
	}
	return ""
}

func (m model) fakeInfoView() string {
	if len(m.fakeInfo) == 0 {
		m.fakeInfo, _ = m.fakeInfoService.GetFakeInfo()
	}
	visibleTopLines := min(max(m.height-3, 0), len(m.fakeInfo))
	fakeBlock := strings.Join(m.fakeInfo[:visibleTopLines], "\n")
	priceLine := fmt.Sprintf("%s %-11s %-21s | %-21s \n%s", "date", "", "Futures", "Actuals", m.futures)
	return fmt.Sprintf("%s\n%s\n[q] quit, [tab] switch", fakeBlock, priceLine)
}

func (m model) historyView() string {
	// futuresHistory view
	history := m.futuresHistory.GetAll()
	visibleLines := max(m.height-6, 0)
	lines := make([]string, visibleLines)
	start := max(len(history)-visibleLines-m.historyScroll, 0)
	end := max(len(history)-m.historyScroll, 0)
	if end > len(history) {
		end = len(history)
	}
	for i, j := start, 0; i < end && j < visibleLines; i, j = i+1, j+1 {
		lines[j] = history[i]
	}
	if len(history) == 0 && visibleLines > 0 {
		lines[0] = "(no history)"
	}
	header := fmt.Sprintf(
		"History in 1 hour (%d entries)\n%s %-11s %-21s | %-21s",
		len(history),
		"date",
		"",
		"Futures",
		"Actuals",
	)
	return fmt.Sprintf("%s\n%s\n%s\n%s\n\n[q] quit, [tab] switch, [up/down] scroll", header, strings.Join(lines, "\n"), "Current:", m.futures)
}

func (m model) chartView() string {
	var sb strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().Bold(true)
	sb.WriteString(titleStyle.Render("Price Trend Chart"))
	sb.WriteString("  ")
	sb.WriteString(m.futures)
	sb.WriteString("\n\n")

	if !m.futuresChartInited && !m.actualsChartInited {
		sb.WriteString("Waiting for data... (first update in ~5 seconds)\n")
	} else {
		// Futures chart
		futuresLabel := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
		sb.WriteString(futuresLabel.Render("▌Futures"))
		if !m.futuresChartInited {
			sb.WriteString("  (market closed)")
		}
		sb.WriteString("\n")
		sb.WriteString(m.futuresChart.View())
		sb.WriteString("\n\n")

		// Actuals chart
		actualsLabel := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
		sb.WriteString(actualsLabel.Render("▌Actuals"))
		if !m.actualsChartInited {
			sb.WriteString("  (no data)")
		}
		sb.WriteString("\n")
		sb.WriteString(m.actualsChart.View())
	}

	sb.WriteString("\n\n[q] quit, [tab] switch")
	return sb.String()
}

// timeLabelFormatter formats X axis labels as HH:MM using local time
func timeLabelFormatter() func(i int, v float64) string {
	return func(i int, v float64) string {
		t := time.Unix(int64(v), 0)
		return t.Format("15:04")
	}
}

// yLabelFormatter formats Y axis labels as integers (price values)
func yLabelFormatter() func(i int, v float64) string {
	return func(i int, v float64) string {
		return fmt.Sprintf("%d", int(v))
	}
}

func newChart(w, h int, lineColor lipgloss.Color) timeserieslinechart.Model {
	chartStyle := lipgloss.NewStyle().Foreground(lineColor)
	axisStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	now := time.Now()

	tslc := timeserieslinechart.New(w, h,
		timeserieslinechart.WithTimeRange(now, now.Add(1*time.Hour)),
		timeserieslinechart.WithAxesStyles(axisStyle, labelStyle),
		timeserieslinechart.WithXYSteps(4, 3),
		timeserieslinechart.WithLineStyle(runes.ArcLineStyle),
		timeserieslinechart.WithStyle(chartStyle),
		timeserieslinechart.WithXLabelFormatter(timeLabelFormatter()),
		timeserieslinechart.WithYLabelFormatter(yLabelFormatter()),
	)
	return tslc
}

func NewProgram(futuresService futures.IService, fakeInfoService fakeinfo.IService) *tea.Program {
	// Get terminal size
	w, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil || height <= 0 {
		height = 24 // Default terminal size
	}
	width := 80
	if err == nil && w > 0 {
		width = w
	}

	fakeInfoService.SetProcessCount(height)

	// Chart dimensions: split vertical space between two charts
	chartWidth := max(width-2, 20)
	chartHeight := max((height-8)/2, 5)

	m := model{
		futuresService:  futuresService,
		fakeInfoService: fakeInfoService,
		fakeInfo:        []string{},
		futures: fmt.Sprintf("[%s] %-21s | %-21s ",
			time.Now().Format("01/02 15:04:05"),
			"-",
			"-",
		),
		height:         height,
		width:          width,
		futuresHistory: shared.NewRingBuffer[string](60),
		futuresChart:   newChart(chartWidth, chartHeight, lipgloss.Color("10")),
		actualsChart:   newChart(chartWidth, chartHeight, lipgloss.Color("12")),
		chartWindowDur: 1 * time.Hour,
	}
	return tea.NewProgram(m, tea.WithAltScreen(), tea.WithInputTTY())
}
