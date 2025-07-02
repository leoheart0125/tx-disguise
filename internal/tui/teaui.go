package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"tx-disguise/internal/futures"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
)

type model struct {
	futuresService futures.IService
	fakeInfo       []string
	futures        string
}

type (
	fakeInfoMsg []string
	futuresMsg  string
)

func genFakeInfoMsg() []string {
	out, err := exec.Command("top", "-stats", "pid,command,cpu,mem", "-l", "2", "-n", "50", "-o", "cpu").Output()
	if err != nil {
		return []string{"[ERROR] " + err.Error()}
	}
	lines := strings.Split(string(out), "\n")
	headerIdx := -1
	header := "Processes:"
	for i, line := range lines {
		if strings.HasPrefix(line, header) {
			if headerIdx == -1 {
				headerIdx = i
				continue
			} else {
				return lines[i-1:]
			}
		}
	}
	return lines
}

func (m model) fakeInfoTicker() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return fakeInfoMsg(genFakeInfoMsg())
	})
}

func (m model) futuresTicker() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return futuresMsg(fmt.Sprintf("[%s] %-21s | %-21s \n",
			time.Now().Format("01/02 15:04:05"),
			m.futuresService.GetCurrentFuturesPrice(),
			m.futuresService.GetCurrentActualPrice(),
		))
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.fakeInfoTicker(), m.futuresTicker())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fakeInfoMsg:
		m.fakeInfo = msg
		return m, m.fakeInfoTicker()

	case futuresMsg:
		m.futures = string(msg)
		return m, m.futuresTicker()

	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	// Get terminal size
	// TODO: width will be used in future updates, default to 80
	_, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil || height <= 0 {
		height = 24 // Default terminal size
	}
	if len(m.fakeInfo) == 0 {
		m.fakeInfo = genFakeInfoMsg()
	}
	visibleTopLines := max(height-3, 0)
	fakeBlock := strings.Join(m.fakeInfo[:visibleTopLines], "\n")
	priceLine := fmt.Sprintf("%s %-11s %-21s | %-21s \n%s", "date", "", "Futures", "Actuals", m.futures)

	return fmt.Sprintf("%s\n%s\n[q] quit", fakeBlock, priceLine)
}

func NewProgram(futuresService futures.IService) *tea.Program {
	m := model{
		futuresService: futuresService,
		fakeInfo:       []string{},
		futures: fmt.Sprintf("[%s] %-21s | %-21s \n",
			time.Now().Format("01/02 15:04:05"),
			"-",
			"-",
		),
	}
	return tea.NewProgram(m, tea.WithAltScreen(), tea.WithInputTTY())
}
