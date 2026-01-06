package wrapper

// Cpu represents CPU information including vendor, model, frequency, and core count.
type Cpu struct {
	Vendor    string
	Model     string
	Frequency string
	MaxCore   uint32
}

// Ram represents RAM information including total and free memory.
type Ram struct {
	Total string
	Free  string
}

// DiskPartition represents a single disk partition with device identifiers and block count.
type DiskPartition struct {
	Name   string
	Major  uint32
	Minor  uint32
	Blocks uint64
}

// Disk represents disk information containing multiple partitions.
type Disk struct {
	Partitions []DiskPartition
}

// Device represents device information including OS version and uptime.
type Device struct {
	OsVersion string
	Uptime    string
}

// ProcessState represents the state of a process.
type ProcessState uint32

const (
	ProcessUndefined      ProcessState = 0
	ProcessRunning        ProcessState = 1
	ProcessSleeping       ProcessState = 2
	ProcessDiskSleep      ProcessState = 3
	ProcessStopped        ProcessState = 4
	ProcessTracingStopped ProcessState = 5
	ProcessZombie         ProcessState = 6
	ProcessDead           ProcessState = 7
)

// String returns a human-readable representation of the process state.
func (ps ProcessState) String() string {
	switch ps {
	case ProcessRunning:
		return "Running"
	case ProcessSleeping:
		return "Sleeping"
	case ProcessDiskSleep:
		return "Disk Sleep"
	case ProcessStopped:
		return "Stopped"
	case ProcessTracingStopped:
		return "Tracing Stopped"
	case ProcessZombie:
		return "Zombie"
	case ProcessDead:
		return "Dead"
	default:
		return "Undefined"
	}
}

// Process represents a single process with its attributes.
type Process struct {
	PID        uint32
	Name       string
	State      ProcessState
	UTime      uint64
	STime      uint64
	NumThreads uint32
}
