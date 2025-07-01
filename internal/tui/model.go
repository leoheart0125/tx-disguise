package tui

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
	"tx-disguise/internal/future"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	futureService future.IService
	fakeInfo      []string
	futures       string
	termHeight    int
}

type (
	fakeInfoMsg []string
	futureMsg   string
)

func (m model) fakeInfoTicker() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		out, err := exec.Command("top", "-status", "pid,command,cpu,mem", "-l", "1", "-n", "50").Output()
		if err != nil {
			return []string{"[ERROR] " + err.Error()}
		}
		lines := strings.Split(string(out), "\n")
		return fakeInfoMsg(lines)
	})
}

func (m model) futureTicker() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return futureMsg(fmt.Sprintf("[%s] %-21s | %-21s \n",
			time.Now().Format("01/02 15:04:05"),
			m.futureService.GetCurrentFuturePrice(),
			m.futureService.GetCurrentActualPrice(),
		))
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.fakeInfoTicker(), m.futureTicker())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		fmt.Println("got WindowSizeMsg", msg.Height)
		m.termHeight = msg.Height

	case fakeInfoMsg:
		m.fakeInfo = msg
		return m, m.fakeInfoTicker()

	case futureMsg:
		m.futures = string(msg)
		return m, m.futureTicker()

	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	// 保留底部 2 行（divider + futures）
	m.termHeight = 30 // Default height for testing, replace with actual terminal height in production
	visibleTopLines := m.termHeight - 3
	if visibleTopLines < 0 {
		visibleTopLines = 0
	}
	if visibleTopLines > len(m.fakeInfo) {
		visibleTopLines = len(m.fakeInfo)
	}

	fakeBlock := strings.Join(m.fakeInfo[:visibleTopLines], "\n")
	priceLine := fmt.Sprintf("%s %-11s %-21s | %-21s \n\n%s", "date", "", "Futures", "Actuals", m.futures)
	divider := strings.Repeat("-", len(priceLine))

	return fmt.Sprintf("%s\n%s\n%s\n[q] quit", fakeBlock, divider, priceLine)
}

func NewProgram(futureService future.IService) *tea.Program {
	m := model{
		futureService: futureService,
		fakeInfo:      []string{},
		futures:       "",
	}
	return tea.NewProgram(m, tea.WithAltScreen(), tea.WithInputTTY())
}
