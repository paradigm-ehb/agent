package svctypes

import (
	"github.com/godbus/dbus"
)

type BusConnection struct {
	conn *dbus.Conn   // pointer to dbus connectoin
	obj  *dbus.Object // pointer to dbus object
}

type System struct {
	os        string    // os version
	processes []Process // list of processes on the server
	services  []Service // list of services on the server
}

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

type Process struct {
	name string // process name
	id   uint32 // PID
}

type Service struct {
	name  string // unit file name
	id    uint32 // service PID
	owner string // unit file owner
}

// type a(ss)
type Ass struct {
	Name  string
	State string
}

// type a(ssssssouso)
type Assssssouso struct {
	Name        string
	Description string
	LoadState   string
	SubState    string
	ActiveState string
	DepUnit     string
	ObjectPath  dbus.ObjectPath
	QueudJob    uint32
	JobType     string
	JobPath     dbus.ObjectPath
}

// type a(sss)
type Asss struct {
	TypeOfChange       string
	FileNameSymLink    string
	DestinationSymLink string
}
