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


