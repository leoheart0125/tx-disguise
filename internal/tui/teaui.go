package tui

import (
	"fmt"
	"os"
	"strings"
	"time"
	"tx-disguise/internal/fakeinfo"
	"tx-disguise/internal/futures"
	"tx-disguise/internal/shared"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
)

type model struct {
	futuresService  futures.IService
	fakeInfoService fakeinfo.IService
	fakeInfo        []string
	futures         string
	futuresHistory  *shared.RingBuffer[string]
	height          int
	currentTab      int // 0: fakeInfo, 1: futuresHistory
	historyScroll   int // scroll position for futuresHistory
}

type (
	fakeInfoMsg       []string
	futuresMsg        string
	futuresHistoryMsg *shared.RingBuffer[string]
)

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
		return futuresMsg(fmt.Sprintf("[%s] %-21s | %-21s ",
			time.Now().Format("01/02 15:04:05"),
			m.futuresService.GetCurrentFuturesPrice(),
			m.futuresService.GetCurrentActualPrice(),
		))
	})
}

func (m model) futuresHistoryTicker() tea.Cmd {
	return tea.Tick(1*time.Minute, func(t time.Time) tea.Msg {
		m.futuresHistory.Push(m.futures)
		return futuresHistoryMsg(m.futuresHistory)
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.fakeInfoTicker(), m.futuresTicker(), m.futuresHistoryTicker())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fakeInfoMsg:
		m.fakeInfo = msg
		return m, m.fakeInfoTicker()

	case futuresMsg:
		m.futures = string(msg)
		return m, m.futuresTicker()

	case futuresHistoryMsg:
		m.futuresHistory = msg
		return m, m.futuresHistoryTicker()

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "tab":
			m.currentTab = (m.currentTab + 1) % 2
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
	if m.currentTab == 0 {
		if len(m.fakeInfo) == 0 {
			m.fakeInfo, _ = m.fakeInfoService.GetFakeInfo()
		}
		visibleTopLines := min(max(m.height-3, 0), len(m.fakeInfo))
		fakeBlock := strings.Join(m.fakeInfo[:visibleTopLines], "\n")
		priceLine := fmt.Sprintf("%s %-11s %-21s | %-21s \n%s", "date", "", "Futures", "Actuals", m.futures)
		return fmt.Sprintf("%s\n%s\n[q] quit, [tab] switch", fakeBlock, priceLine)
	}
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

func NewProgram(futuresService futures.IService, fakeInfoService fakeinfo.IService) *tea.Program {
	// Get terminal size
	// TODO: width will be used in future updates, default to 80
	_, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil || height <= 0 {
		height = 24 // Default terminal size
	}
	fakeInfoService.SetProcessCount(height)
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
		futuresHistory: shared.NewRingBuffer[string](60),
	}
	return tea.NewProgram(m, tea.WithAltScreen(), tea.WithInputTTY())
}
