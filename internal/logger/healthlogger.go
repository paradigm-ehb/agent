package logger

// TODO(nasr): cflags to link to the binary
/*
* include <stdlib.h>
* include agent-resources.h>
* */

import (
	"C"
	"fmt"
)

type LogLevel int

type ServerState int

// enum server healtha
// @param
// Healthy
// Unhealthy
// AttentionNeeded
const (
	Healthy ServerState = iota
	Unhealthy
	AttentionNeeded
)

const (
	High LogLevel = iota
	Medium
	Low
	Stable
)

// returns the delta of the cpu frquency
func compareCpuFrequency(cpuFreqX float32, cpuFreqY float32) float32 {

	return cpuFreqY / cpuFreqX
}

func compareCpuUsage(cpuUsageX float32, cpuUsageY float32) float32 {

	return cpuUsageY / cpuUsageX
}

func Status(freqDelta float32, usageDelta float32) LogLevel {

	fmt.Println("checking if cpu is stable")
	if freqDelta < 0 && usageDelta < 0 {
		return Stable
	}

	return Medium
}

func Construct() string {

	// TODO(nasr): construct logs based on passed status arguments
	return "The server is stable and running"
}

// Append the logs to a configuration file
func Append() {

	// TODO(nasr): create logs and write them to the system
}
