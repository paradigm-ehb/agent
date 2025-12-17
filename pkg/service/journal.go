package service

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
	// TODO(nasr): remove the magic number enums, horrible code practice
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

	// Generated time. testing issue
	// TODO(nasr): fix the time
	//sinceTime, err := time.Parse(time.RFC3339, "0")
	//
	//if err != nil {
	//	fmt.Println("error parsing time:", err)
	//	return nil, nil
	//}
	//
	//duration := time.Since(sinceTime)

	var wg sync.WaitGroup

	ch := make(chan []byte)

	wg.Add(2)
	go func() {
		defer wg.Done()
		j.GetJournalInformation(0, in.NumFromTail, in.Cursor, m, in.Path, ch)
	}()

	go func() {
		defer wg.Done()

		for range ch {
			resp := journal.JournalChunk{Reply: <-ch}
			if err := srv.Send(&resp); err != nil {
				log.Printf("send error %v", err)
			}
		}
	}()

	wg.Wait()
	return nil
}
