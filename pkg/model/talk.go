package model

import "time"

type TalkType int

const (
	TalkType_OnlineSession TalkType = iota
	TalkType_RecordingSession
	TalkType_Opening
	TalkType_Ending
	TalkType_Commercial
)

type Talks []Talk

func (ts Talks) WillStartNextTalkSince(d time.Duration) bool {
	// TODO (#10)
	return false
}

func (ts Talks) GetNextTalk() Talk {
	// TODO (#10)
	return Talk{}
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
