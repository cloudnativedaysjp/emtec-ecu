package model

import (
	"time"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/utils"
)

type TalkType int32

const (
	TalkType_OnlineSession TalkType = iota
	TalkType_RecordingSession
	TalkType_Opening
	TalkType_Ending
	TalkType_Commercial
)

type Talks []Talk

func (ts *Talks) FillCommercial() {
	// TODO (#10)
}

func (ts Talks) WillStartNextTalkSince(d time.Duration) bool {
	// TODO (#10)
	return false
}

func (ts Talks) GetCurrentTalk() Talk {
	// TODO (#10)
	return Talk{}
}

func (ts Talks) GetNextTalk() Talk {
	// TODO (#10)
	return Talk{}
}

func (ts Talk) GetActualStartAtAndEndAt(conferenceDayDate string, actualStartTime, actualEndTime time.Time) (*time.Time, *time.Time, error) {
	cDate, err := utils.ParseDateFormat(conferenceDayDate)
	if err != nil {
		return nil, nil, err
	}

	return time.Time, time.Time
}

type Talk struct {
	Id           int32
	TalkName     string
	TrackId      int32
	TrackName    string
	EventAbbr    string
	SpeakerNames []string
	Type         TalkType
	StartAt      time.Time
	EndAt        time.Time
}
