package service

import (
	"context"
	"fmt"
	"paradigm-ehb/agent/gen/journal/v1"
	j "paradigm-ehb/agent/internal/journal"

	sdjournal "github.com/coreos/go-systemd/v22/sdjournal"
)

type JournalService struct {
	journal.UnimplementedJournalServiceServer
}

func (s *JournalService) Action(_ context.Context, in *journal.JournalRequest) (*journal.JournalReply, error) {

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

	output, err := j.GetJournalInformation(0, in.NumFromTail, in.Cursor, m, in.Path)
	if err != nil {
		fmt.Println("error occurred when calling GetJournalInformation:", err)
	}
	return &journal.JournalReply{Reply: output}, nil
}
