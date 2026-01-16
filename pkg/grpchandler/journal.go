package grpc_handler

import (
	"log"
	"paradigm-ehb/agent/gen/journal/v1"
	j "paradigm-ehb/agent/internal/journal"
	"sync"

	"github.com/coreos/go-systemd/v22/sdjournal"
)

type JournalService struct {
	journal.UnimplementedJournalServiceServer
}

func (s *JournalService) Action(in *journal.JournalRequest, srv journal.JournalService_ActionServer) error {

	var val string
	/*
	TODO(nasr):
	- Remove magic-number–based enum handling.
	- Replace `switch in.Field` with:
	  - a typed protobuf enum, or
	  - a map[JournalField]string lookup, or
	  - explicit constants with semantic names.
	- Ensure invalid enum values are handled (default case + error).
	*/
	switch in.Field {

	case 0:
		val = sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT
	case 1:
		val = sdjournal.SD_JOURNAL_FIELD_PID
	case 2:
		val = sdjournal.SD_JOURNAL_FIELD_UID
	case 3:
		val = sdjournal.SD_JOURNAL_FIELD_GID
	}

	m := []sdjournal.Match{{Field: val, Value: in.Value}}

	/*
	TODO(nasr):
	- Implement proper time filtering support.
	- Decide on API semantics:
	  - absolute timestamp vs relative duration
	  - server-side vs client-provided time window
	- Validate time parsing errors and propagate them via gRPC status.
	- Remove dead/commented-out code once design is finalized.
	*/

	var wg sync.WaitGroup

	ch := make(chan []byte)

	/*
	TODO(nasr):
	- Define clear channel ownership:
	  - who closes `ch` and when.
	- Consider making channel buffered to avoid producer/consumer blocking.
	- Document lifetime guarantees between goroutines.
	*/

	wg.Add(2)
	go func() {
		defer wg.Done()
		/*
		TODO(nasr):
		- Propagate context cancellation into GetJournalInformation.
		- Return errors instead of silent failure.
		- Clarify meaning of leading `0` argument.
		*/
		j.GetJournalInformation(0, in.NumFromTail, in.Cursor, m, in.Path, ch)
	}()

	go func() {
		defer wg.Done()

		/*
		TODO(nasr):
		- Fix channel consumption logic:
		  current code ranges over `ch` but also performs `<-ch` again.
		- Avoid double reads from the same channel.
		- Handle channel close explicitly.
		- Propagate send errors via context cancellation or status return.
		*/
		for range ch {
			resp := journal.JournalChunk{Reply: <-ch}
			if err := srv.Send(&resp); err != nil {
				log.Printf("send error %v", err)
			}
		}
	}()

	/*
	TODO(nasr):
	- Consider early exit on client disconnect.
	- Avoid waiting indefinitely if goroutines deadlock.
	- Evaluate replacing WaitGroup + channels with errgroup + context.
	*/
	wg.Wait()
	return nil
}
