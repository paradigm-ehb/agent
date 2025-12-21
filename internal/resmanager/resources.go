package resources

/*
#cgo CFLAGS:  -I${SRCDIR}/agent-resources
#cgo LDFLAGS: -L${SRCDIR}/agent-resources/build -lagent_resources
#include "agent_resources.h"
*/
import "C"
import (
	"log"
)

// TODO(nasr): add fields
type CPUArena struct {
	Vendor    string
	Model     string
	Frequency string
	Cores     string
}
type MemoryArena struct {
	Total string
	Free  string
}
type DiskArena struct{}

type DeviceArena struct {
	OsVersion      string
	Uptime         string
	Processes      []string
	ProcessesCount int
}

/**
* Resource data objects
* */
type CPUUtilization float32
type MemoryUtilization float32
type DiskUtilization float32
type DeviceUtilization float32

/**
*
* Create a snapshot of the current CPU and RAM state
* to determine if the system is stable
*
* */
type SystemSnapshot struct {
	CPUFrequency CPUUtilization
	MemoryUsage  MemoryUtilization
	Timestamp    int64
}

/**
*
* HealthLevel assigns a value to the determined level of stability
* */
type HealthLevel int

const (
	Stable HealthLevel = iota
	Warning
	Critical
)

/*
* Check's for CPU drops.
* An overheating or non-stable CPU drops clock frequency.
* This will determine if the CPU is stable under the current load
* CalculateCPUFrequencyRatio computes the ratio between two CPU frequency snapshots
* @return: ratio where >1 means frequency increased, <1 means decreased
*
**/

func CreateAgentRam() {

	ram := C.agent_ram_create()
	if ram == nil {
		return
	}
	defer C.agent_ram_destroy(ram)

	if C.agent_ram_read(ram) != C.AGENT_OK {
		return
	}

	log.Println("total ram -> ", C.GoString(C.agent_ram_get_total(ram)))
	log.Println("free ram -> ", C.GoString(C.agent_ram_get_free(ram)))
}

func CreateAgentCpu() {

	cpu := C.agent_cpu_create()

	if cpu == nil {
		return
	}

	defer C.agent_cpu_destroy(cpu)

	if C.agent_cpu_read(cpu) != C.AGENT_OK {
		return
	}

	// TODO(nasr): currently not available on arm
	log.Println("cpu -> ", C.GoString(C.agent_cpu_get_vendor(cpu)))
	log.Println("cpu -> ", C.GoString(C.agent_cpu_get_model(cpu)))
	log.Println("cpu -> ", C.GoString(C.agent_cpu_get_frequency(cpu)))
	log.Println("cpu -> ", C.GoString(C.agent_cpu_get_cores(cpu)))

}

func CreateAgentDevice() {

	device := C.agent_device_create()

	if device == nil {
		return
	}

	defer C.agent_device_destroy(device)

	if C.agent_device_read(device) != C.AGENT_OK {
		return

	}

	log.Println("device ->", C.GoString(C.agent_device_get_uptime(device)))
	log.Println("device -> ", C.GoString(C.agent_device_get_os_version(device)))

}

func CalculateCPUFrequencyRatio(baseline, current SystemSnapshot) float32 {
	if baseline.CPUFrequency == 0 {
		return 0 // Prevent division by zero
	}

	return float32(current.CPUFrequency) / float32(baseline.CPUFrequency)
}

// CalculateCPUUsageRatio computes the ratio between two CPU usage measurements
// Returns: ratio where >1 means usage increased, <1 means decreased
func CalculateCPUUsageRatio(baseline, current float32) float32 {
	if baseline == 0 {
		return 0 // Prevent division by zero
	}
	return current / baseline
}

/**
 * CaptureSystemSnapshot retrieves current system resource state
 * using the agent-resources libarary
 */
func CaptureSystemSnapshot() (SystemSnapshot, error) {
	return SystemSnapshot{}, nil
}

// EvaluateSystemHealth determines system stability based on resource deltas
// frequencyRatio: ratio of current/baseline CPU frequency
// usageRatio: ratio of current/baseline CPU usage
func EvaluateSystemHealth(frequencyRatio, usageRatio float32) HealthLevel {
	log.Println("evaluating system stability")

	if frequencyRatio < 1.0 && usageRatio < 1.0 {
		return Critical
	}

	if frequencyRatio < 0.8 || usageRatio > 1.5 {
		return Warning
	}

	return Stable
}
