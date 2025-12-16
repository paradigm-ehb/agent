package logger

// TODO(nasr): cflags to link to the binary
/*
* include <stdlib.h>
* include agent-resources.h>
* */

import (
	// "C"
	"fmt"
	jrnl "github.com/coreos/go-systemd/v22/journal"
	sdj "github.com/coreos/go-systemd/v22/sdjournal"
)

type LogLevel int

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

type ServerState int

const (
	High LogLevel = iota
	Medium
	Low
	Stable
)

// TODO(nasr): return the current cpu state
type CpuState int

// TODO(nasr): return the current memory state
type MemoryState int

// returns the delta of the cpu frquency
// used to check if the cpu is stable
func compareCpuFrequency(cpuFreqX float32, cpuFreqY float32) float32 {

	// TODO(nasr): retrieve the data with CGO
	return cpuFreqY / cpuFreqX
}

// returns the deleta of the cpu usage
// used to check if the cpu is stable
func compareCpuUsage(cpuUsageX float32, cpuUsageY float32) float32 {

	// TODO(nasr): retrieve the data with CGO
	return cpuUsageY / cpuUsageX
}

// TODO(nasr): call cgo to create a snapshot of the current system
func CreateSnapshot() {

}

func Status(freqDelta float32, usageDelta float32) LogLevel {

	fmt.Println("checking if cpu is stable")
	if freqDelta < 0 && usageDelta < 0 {
		return Stable
	}

	return Medium
}

// Function to check if systemd is enabled on the device
func CheckJournal() bool {

	return jrnl.Enabled()
}

// Append the logs to a configuration file
func foo() {

	j, err := sdj.NewJournal()
	if err != nil {
		fmt.Printf("%v->", j)
	}

	bid := j.GetBootID
	fmt.Println("output: ", bid)

	err = j.Close()
	if err != nil {
		fmt.Println("failed to close the journal")
	}
}

func journalConfig() *sdj.JournalReaderConfig {

	return &sdj.JournalReaderConfig{Since: 0, NumFromTail: 10, Matches: []sdj.Match{{Field: "_SYSTEMD_UNIT=", Value: "nginx.service"}}, Path: ""}
}

func find(config *sdj.JournalReaderConfig) {

	reader, err := sdj.NewJournalReader(*config)
	if err != nil {
		fmt.Println("failed to make the reader")
	}

	var b []byte
	code, err := reader.Read(b)
	if err != nil && code != 0 {
		fmt.Println("Hello World")
	}
}

func Run() {

	config := journalConfig()
	find(config)

}
