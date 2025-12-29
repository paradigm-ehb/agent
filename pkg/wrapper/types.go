package wrapper

type Cpu struct {
	Vendor    string
	Model     string
	Frequency string
	MaxCore   uint32
}

type Ram struct {
	Total string
	Free  string
}

type DiskPartition struct {
	Name   string
	Major  uint32
	Minor  uint32
	Blocks uint64
}

type Disk struct {
	Partitions []DiskPartition
}

type Device struct {
	OsVersion string
	Uptime    string
}

type ProcessState rune

const (
	Running        ProcessState = 'R'
	Sleeping       ProcessState = 'S'
	DiskSleeping   ProcessState = 'D'
	Stopped        ProcessState = 'T'
	TracingStopped ProcessState = 't'
	Zombie         ProcessState = 'Z'
	Dead           ProcessState = 'X'
)

type Process struct {
	PID        uint32
	Name       string
	State      ProcessState
	UTime      uint64
	NumThreads uint32
}
