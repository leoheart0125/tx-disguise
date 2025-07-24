package fakeinfo

import (
	"os/exec"
	"strings"
)

type TopFakeInfoService struct {
	command      []string
	processCount int
}

func NewTopFakeInfoService(command string) *TopFakeInfoService {
	return &TopFakeInfoService{
		command: []string{"top", "-stats", "pid,command,cpu,mem", "-l", "2", "-n", "50", "-o", "cpu"},
	}
}

var _ IService = (*TopFakeInfoService)(nil)

func (s *TopFakeInfoService) SetProcessCount(count int) {
	if count < 1 {
		count = 24
	}
	s.processCount = count
}

func (s *TopFakeInfoService) GetFakeInfo() ([]string, error) {
	out, err := exec.Command(s.command[0], s.command[1:]...).Output()
	if err != nil {
		return []string{"[ERROR] " + err.Error()}, err
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
				return lines[i:s.processCount], nil
			}
		}
	}
	return lines, nil
}
