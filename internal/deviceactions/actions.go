package deviceactions

type DeviceAction string

const (
	DeviceActionShutdown  DeviceAction = "org.freedesktop.login1.Manager.PowerOff"  // shutdown device
	DeviceActionReboot    DeviceAction = "org.freedesktop.login1.Manager.Reboot"    // reboot device
	DeviceActionSuspend   DeviceAction = "org.freedesktop.login1.Manager.Suspend"   // suspend device
	DeviceActionHibernate DeviceAction = "org.freedesktop.login1.Manager.Hibernate" // hibernate device
)
