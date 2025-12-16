package logger

import (
	"fmt"
	jrnl "github.com/coreos/go-systemd/v22/journal"
	sdj "github.com/coreos/go-systemd/v22/sdjournal"
	"time"
)

// checkJournal reports whether systemd’s journal is available and enabled
// on the current system.
//
// It is a lightweight capability check and does not open or read the journal.
// Internally, this relies on libsystemd to detect whether journald is usable
// (for example, not present on non-systemd systems).
func checkJournal() bool {

	return jrnl.Enabled()
}

// systemdID returns the boot ID associated with the currently running system.
//
// The boot ID uniquely identifies the current boot session and is useful for
// correlating journal entries to a specific system start. The function opens
// the journal, queries the boot ID, and then closes the journal handle.
//
// If the journal cannot be opened or the boot ID cannot be retrieved, the
// returned string may be empty.
func systemdID() string {

	j, err := sdj.NewJournal()
	if err != nil {
		fmt.Printf("%v->", j)
	}

	bid, err := j.GetBootID()
	if err != nil {
		fmt.Println("failed to get the bootd id")
	}

	err = j.Close()
	if err != nil {
		fmt.Println("failed to close the journal")
	}

	return bid

}

// TODO(nasr): checkout formatters

// GetJournaldInformation reads entries from the systemd journal and returns
// them as a single concatenated string.
//
// The journal reader is configured through the provided parameters:
//   - since:        limits entries to those newer than the given duration
//     relative to now
//   - numFromTail:  limits the number of entries read from the end of the journal
//   - cursor:       reserved for future cursor-based positioning (currently unused)
//   - matches:      filters entries using systemd journal match rules
//   - path:         optionally specifies a custom journal path
//
// Internally, this function uses a JournalReader and performs sequential reads
// into a fixed-size buffer until no more data is available or an error occurs.
// The caller receives raw journal output as text, without further parsing or
// field-level decoding.
//
//	Example Matches:     []sdj.Match{{Field: "_SYSTEMD_UNIT", Value: "ssh.service"}}}
func GetJournaldInformation(since time.Duration, numFromTail uint64, cursor string, matches []sdj.Match, path string) string {

	config := sdj.JournalReaderConfig{
		Since:       since,
		NumFromTail: numFromTail,
		Matches:     matches,
		Path:        path,
	}

	reader, err := sdj.NewJournalReader(config)

	if err != nil {
		fmt.Println("Failed to open the journalctl reader")
	}

	defer reader.Close()

	b := make([]byte, 4096)
	var output string

	for {
		c, err := reader.Read(b)
		if err != nil {
			fmt.Printf("\nfailed when reading, %v ", err)
			break
		}
		if c == 0 {
			continue
		}

		output += string(b[:])
	}

	return output
}
