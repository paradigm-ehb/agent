package service

import (
	"context"
	"fmt"

	sdj "github.com/coreos/go-systemd/v22/sdjournal"
	"paradigm-ehb/agent/gen/journal/v1"
	j "paradigm-ehb/agent/internal/journal"
	t "paradigm-ehb/agent/internal/journal/types"
	"time"
)

type JournalService struct {
	journal.UnimplementedJournalServiceServer
}

func (s *JournalService) Action(_ context.Context, in *journal.JournalRequest) (*journal.JournalReply, error) {

	var val string

	fmt.Println(string(t.Systemd))
	switch in.Field {

	case 0:
		val = string(t.Systemd)
	case 1:
		val = string(t.PID)
	case 2:
		val = string(t.UID)
	case 3:
		val = string(t.GID)
		break

	}

	m := []sdj.Match{
		{
			Field: val,
			Value: in.Value,
		},
	}

	fmt.Println(m)
	duration, err := time.ParseDuration(in.Time)
	if err != nil {
		return nil, nil
	}
	log := j.GetJournalInformation(duration, in.NumFromTail, in.Cursor, m, in.Path)
	fmt.Println("log: ", log)
	return &journal.JournalReply{Reply: log}, nil
}
