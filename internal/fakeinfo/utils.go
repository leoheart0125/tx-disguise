package fakeinfo

import "math/rand"

func CPUUsageGenerator() float64 {
	// This function simulates CPU usage generation.
	weight := rand.Float64()
	max, min := 0, 0
	switch {
	case weight < 0.01:
		max, min = 100, 80 // High CPU usage
	case weight < 0.4:
		max, min = 30, 5 // Moderate CPU usage
	default:
		max, min = 10, 0 // Low CPU usage
	}
	return rand.Float64()*float64(max-min) + float64(min)
}

func MemoryUsageGenerator() int {
	// This function simulates memory usage generation.
	weight := rand.Float64()
	switch {
	case weight < 0.1:
		return rand.Intn(1000) + 1000 // High memory usage between 1000 and 1999
	case weight < 0.5:
		return rand.Intn(500) + 100 // Moderate memory usage between 100 and 599
	default:
		return rand.Intn(100) + 100 // Low memory usage between 100 and 199
	}
}
