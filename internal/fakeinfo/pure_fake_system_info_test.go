package fakeinfo

import (
	"strings"
	"testing"
)

func TestNewPureFakeSystemInfoService_DefaultProcessCount(t *testing.T) {
	svc := NewPureFakeSystemInfoService()
	if svc.processCount != 24 {
		t.Errorf("Expected default processCount 24, got %d", svc.processCount)
	}
}

func TestSetProcessCount(t *testing.T) {
	svc := NewPureFakeSystemInfoService()
	svc.SetProcessCount(10)
	if svc.processCount != 10 {
		t.Errorf("Expected processCount 10, got %d", svc.processCount)
	}

	svc.SetProcessCount(0)
	if svc.processCount != 24 {
		t.Errorf("Expected fallback processCount 24, got %d", svc.processCount)
	}
}

func TestGenerateSystemInfo_Format(t *testing.T) {
	svc := NewPureFakeSystemInfoService()
	info := svc.generateSystemInfo()
	if !strings.Contains(info, "Processes:") {
		t.Error("System info should contain 'Processes:'")
	}
	if !strings.Contains(info, "CPU usage:") {
		t.Error("System info should contain 'CPU usage:'")
	}
}

func TestGenerateFakeProcesses_Count(t *testing.T) {
	svc := NewPureFakeSystemInfoService()
	procs := svc.generateFakeProcesses(5)
	if len(procs) != 5 {
		t.Errorf("Expected 5 processes, got %d", len(procs))
	}
	for _, p := range procs {
		if p.PID < 1 {
			t.Error("PID should be >= 1")
		}
		if p.Memory < 100 {
			t.Error("Memory should be >= 100")
		}
	}
}

func TestGetFakeInfo_Output(t *testing.T) {
	svc := NewPureFakeSystemInfoService()
	svc.SetProcessCount(5)
	info, err := svc.GetFakeInfo()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(info) == 0 {
		t.Error("Expected non-empty info output")
	}
	foundHeader := false
	for _, line := range info {
		if strings.Contains(line, "PID") && strings.Contains(line, "COMMAND") {
			foundHeader = true
			break
		}
	}
	if !foundHeader {
		t.Error("Expected process table header in output")
	}
}
