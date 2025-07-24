package fakeinfo

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type FakeProcess struct {
	PID     int
	Command string
	CPU     float64
	Memory  int
}

type PureFakeSystemInfoService struct {
	processCount int
}

var _ IService = (*PureFakeSystemInfoService)(nil)

func NewPureFakeSystemInfoService() *PureFakeSystemInfoService {
	return &PureFakeSystemInfoService{
		processCount: 24,
	}
}

func (s *PureFakeSystemInfoService) SetProcessCount(count int) {
	if count < 1 {
		count = 24
	}
	s.processCount = count
}

func (s *PureFakeSystemInfoService) generateSystemInfo() string {
	const systemInfo = `
Processes: %d total, %d running, %d sleeping, %d threads
%s
Load Avg: %.2f, %.2f, %.2f
CPU usage: %.2f%% user, %.2f%% sys, %.2f%% idle
SharedLibs: %dM resident, %dM data, %dM linkedit.
MemRegions: %d total, %dB resident, %dM private, %dM shared.
PhysMem: %dG used (%dM wired, %dM compressor), %dM unused.
VM: %dT vsize, %dM framework vsize, %d(%d) swapins, %d(%d) swapouts.
Networks: packets: %d/%dM in, %d/%dM out.
Disks: %d/%dG read, %d/%dG written.
`
	totalProcs := rand.Intn(300) + 100
	runningProcs := rand.Intn(totalProcs/2) + 10
	sleepingProcs := totalProcs - runningProcs
	threads := rand.Intn(1000) + 500

	loadAvg1 := rand.Float64() * 5
	loadAvg2 := rand.Float64() * 5
	loadAvg3 := rand.Float64() * 5

	cpuUser := rand.Float64() * 50
	cpuSys := rand.Float64() * 20
	cpuIdle := 100 - cpuUser - cpuSys

	sharedLibsResident := rand.Intn(1000)
	sharedLibsData := rand.Intn(200)
	sharedLibsLinkedit := rand.Intn(300)

	memRegionsTotal := rand.Intn(10)
	memRegionsResident := rand.Intn(1000)
	memRegionsPrivate := rand.Intn(1000) + 100
	memRegionsShared := rand.Intn(4000) + 1000

	physMemUsed := rand.Intn(32) + 8
	physMemWired := rand.Intn(4000)
	physMemCompressor := rand.Intn(8000)
	physMemUnused := rand.Intn(4000)

	vmVsize := rand.Intn(500)
	vmFrameworkVsize := rand.Intn(10000)
	swapins := rand.Intn(10)
	swapins2 := rand.Intn(10)
	swapouts := rand.Intn(10)
	swapouts2 := rand.Intn(10)

	netPacketsIn := rand.Intn(50000000)
	netPacketsInM := rand.Intn(10000)
	netPacketsOut := rand.Intn(50000000)
	netPacketsOutM := rand.Intn(10000)

	diskRead := rand.Intn(10000000)
	diskReadG := rand.Intn(200)
	diskWritten := rand.Intn(10000000)
	diskWrittenG := rand.Intn(200)

	return fmt.Sprintf(systemInfo,
		totalProcs, runningProcs, sleepingProcs, threads,
		time.Now().Format("2006/01/02 15:04:05"),
		loadAvg1, loadAvg2, loadAvg3,
		cpuUser, cpuSys, cpuIdle,
		sharedLibsResident, sharedLibsData, sharedLibsLinkedit,
		memRegionsTotal, memRegionsResident, memRegionsPrivate, memRegionsShared,
		physMemUsed, physMemWired, physMemCompressor, physMemUnused,
		vmVsize, vmFrameworkVsize, swapins, swapins2, swapouts, swapouts2,
		netPacketsIn, netPacketsInM, netPacketsOut, netPacketsOutM,
		diskRead, diskReadG, diskWritten, diskWrittenG,
	)
}

func (s *PureFakeSystemInfoService) generateFakeProcesses(count int) []FakeProcess {
	processes := make([]FakeProcess, count)
	for i := range count {
		processes[i] = FakeProcess{
			PID:     rand.Intn(10000) + 1,
			Command: strings.ToLower(gofakeit.AppName()),
			CPU:     rand.Float64() * 30,
			Memory:  rand.Intn(1000) + 100,
		}
	}
	return processes
}

func (s *PureFakeSystemInfoService) GetFakeInfo() ([]string, error) {
	// Generate a random system info string
	var fakeInfo []string
	fakeInfo = append(fakeInfo, strings.Split(s.generateSystemInfo(), "\n")...)
	fakeInfo = append(fakeInfo, fmt.Sprintf("%-6s %-22s %-21s %-21s", "PID", "COMMAND", "%CPU", "MEM"))
	for _, proc := range s.generateFakeProcesses(s.processCount) {
		cmd := proc.Command
		if len(cmd) > 20 {
			cmd = cmd[:20]
		}
		fakeInfo = append(fakeInfo, fmt.Sprintf("%-6d %-22s %-21.1f %-21d", proc.PID, cmd, proc.CPU, proc.Memory))
	}

	return fakeInfo, nil
}
