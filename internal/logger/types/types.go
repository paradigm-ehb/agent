package loggertypes

// Trusted journald fields
type TrustedJournaldField string

const (
	Systemd TrustedJournaldField = "_SYSTEMD_UNIT"
	PID     TrustedJournaldField = "_PID"
	UID     TrustedJournaldField = "_UID"
	GID     TrustedJournaldField = "_GID"
)

type Match struct {
	Field TrustedJournaldField
	Value string
}
