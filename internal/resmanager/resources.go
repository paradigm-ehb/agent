package resources

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

// func Status(freqDelta float32, usageDelta float32) LogLevel {
//
// 	log.Println("checking if cpu is stable")
// 	if freqDelta < 0 && usageDelta < 0 {
// 		return
// 	}
//
// 	return Medium
// }
