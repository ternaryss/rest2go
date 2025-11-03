package rest2go

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

var pidStart = time.Now()

type health struct {
	Timestamp       string  `json:"timestamp"`
	Version         string  `json:"version"`
	Uptime          float64 `json:"uptime"`
	Cpu             float64 `json:"cpu"`
	MemoryUsed      float64 `json:"memoryUsed"`
	MemoryAvailable float64 `json:"memoryAvailable"`
	MemoryTotal     float64 `json:"memoryTotal"`
	Routines        int     `json:"routines"`
}

func newHealth(uptime, cpu, memUsed, memAvailable, memTotal float64, routines int, version string) health {
	return health{
		Timestamp:       time.Now().Format(time.RFC3339),
		Version:         version,
		Uptime:          uptime,
		Cpu:             cpu,
		MemoryUsed:      memUsed,
		MemoryAvailable: memAvailable,
		MemoryTotal:     memTotal,
		Routines:        routines,
	}
}

func toMb(value uint64) float64 {
	return float64(value) / 1024.0 / 1024.0
}

func healthCheck(response http.ResponseWriter, request *http.Request) {
	cpuPercent, memUsed, memAvailable, memTotal := 0.0, 0.0, 0.0, 0.0
	ctx, cancel := context.WithTimeout(request.Context(), 500*time.Millisecond)
	defer cancel()

	if cpuUsage, err := cpu.PercentWithContext(ctx, 0, false); err == nil {
		if len(cpuUsage) > 0 {
			cpuPercent = cpuUsage[0]
		}
	} else {
		HandleError(err, response)
		return
	}

	if memory, err := mem.VirtualMemoryWithContext(ctx); err == nil {
		memUsed = toMb(memory.Used)
		memAvailable = toMb(memory.Available)
		memTotal = toMb(memory.Total)
	} else {
		HandleError(err, response)
		return
	}

	stats := newHealth(
		time.Since(pidStart).Seconds(),
		cpuPercent,
		memUsed,
		memAvailable,
		memTotal,
		runtime.NumGoroutine(),
		runtime.Version(),
	)
	jsonBytes, err := json.Marshal(stats)

	if err != nil {
		HandleError(err, response)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonBytes)
}
