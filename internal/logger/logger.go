package logger

// TODO(nasr): cflags to link to the binary
/*
* include <stdlib.h>
* include agent-resources.h>
* */

import (
	// "C"
	jrnl "github.com/coreos/go-systemd/v22/journal"
	sdj "github.com/coreos/go-systemd/v22/sdjournal"
	"log"
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

	log.Println("checking if cpu is stable")
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
		log.Printf("%v->", j)
	}

	bid := j.GetBootID
	log.Println("output: ", bid)

	err = j.Close()
	if err != nil {
		log.Println("failed to close the journal")
	}
}

func Run() {

	config := sdj.JournalReaderConfig{
		NumFromTail: 10,
		Matches:     []sdj.Match{{Field: "_SYSTEMD_UNIT", Value: "ssh.service"}}}

	reader, err := sdj.NewJournalReader(config)

	if err != nil {
		log.Println("failed to make the reader")
	}
	defer reader.Close()

	b := make([]byte, 4096)

	for {
		c, err := reader.Read(b)
		if err != nil {
			log.Printf("\nfailed when reading, %v ", err)
			break
		}
		if c == 0 {
			continue
		}

		log.Println("here: ", string(b[:]))
		log.Println("======================")
	}
}
