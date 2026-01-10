package journal

import (
	"io"
	"time"

	"log"

	"github.com/coreos/go-systemd/v22/journal"
	"github.com/coreos/go-systemd/v22/sdjournal"
)

/*
checkJournal reports whether systemd’s journal is available and enabled
on the current system.

It is a lightweight capability check and does not open or read the journal.
Internally, this relies on libsystemd to detect whether journald is usable
(for example, not present on non-systemd systems).
*/
func checkJournal() bool {

	return journal.Enabled()
}

/*
systemdID returns the boot ID associated with the currently running system.

The boot ID uniquely identifies the current boot session and is useful for
correlating journal entries to a specific system start. The function opens
the journal, queries the boot ID, and then closes the journal handle.

If the journal cannot be opened or the boot ID cannot be retrieved, the
returned string may be empty.
*/
func systemdID() (string, error) {

	j, err := sdjournal.NewJournal()
	if err != nil {
		return "not available", err
	}

	bid, err := j.GetBootID()
	if err != nil {
		j.Close()
		return "not available", err
	}

	return bid, nil
}

/*
TODO(nasr):
- Evaluate and implement journal output formatters.
- Decide whether formatting should be:
  - raw text (current behavior),
  - structured (map / proto),
  - or selectable via configuration.
*/

/*
GetJournalInformation reads entries from the systemd journal and streams them
through the provided output channel.

Parameters:
- since:
  Duration relative to "now" used to limit journal entries.
  (Currently not implemented.)
- numFromTail:
  Limits the number of entries read from the end of the journal.
- cursor:
  Reserved for cursor-based positioning (currently unused).
- matches:
  Journal match filters (field/value pairs).
- path:
  Optional custom journal path.
- out:
  Output channel receiving raw journal byte slices.

Behavior:
- Opens a JournalReader with the provided configuration.
- Reads sequentially into a fixed-size buffer.
- Streams raw journal output without parsing or decoding fields.
- Closes the output channel before returning.

Example:
  []sdj.Match{{Field: "_SYSTEMD_UNIT", Value: "ssh.service"}}
*/
func GetJournalInformation(
	since time.Duration,
	numFromTail uint64,
	cursor string,
	matches []sdjournal.Match,
	path string,
	out chan []byte,
) {

	/*
	TODO(nasr):
	- Define ownership and lifecycle rules for `out`.
	- Document that this function is responsible for closing the channel.
	*/
	defer close(out)

	config := sdjournal.JournalReaderConfig{
		/*
		TODO(nasr):
		- Implement time-based filtering using `since`.
		- Decide whether `since` should override cursor semantics.
		*/
		// Since:       since,
		NumFromTail: numFromTail,
		Cursor:      cursor,
		Matches:     matches,
		Path:        path,
		Formatter:   nil,
	}

	reader, err := sdjournal.NewJournalReader(config)
	if err != nil {
		/*
		TODO(nasr):
		- Replace stdout logging with structured error propagation.
		- Decide whether to terminate early or send error markers via channel.
		*/

	}

	defer reader.Close()

	b := make([]byte, 4096)

	for {
		c, err := reader.Read(b)

		if err == io.EOF {
			/*
			TODO(nasr):
			- Decide whether EOF should be silent.
			- Avoid stdout logging in library code.
			*/
			log.Printf("end of journal %v, ", err)
			break
		}

		if c == 0 {
			/*
			TODO(nasr):
			- Clarify whether zero-length reads are expected.
			- Remove noisy logging or replace with debug-level tracing.
			*/
			continue
		}

		if err != nil {
			/*
			TODO(nasr):
			- Avoid sending placeholder data on error.
			- Define a structured error signaling mechanism.
			- Consider context cancellation or error channels.
			*/
			out <- []byte("nothing in here")
			log.Printf("no more data to read %v, ", err)
			break
		}

		out <- b[:c]
	}
}
