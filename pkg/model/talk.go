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

func (ts Talks) WillStartNextTalkSince() bool {
	now := time.Now()
	for _, talk := range ts {
		if now.After(talk.StartAt) && talk.EndAt.Sub(now) <= utils.HowManyMinutesUntilNotify {
			return true
		}
	}
	return false
}

func (ts Talks) GetCurrentTalk() (*Talk, int) {
	now := time.Now()
	for i, talk := range ts {
		if now.After(talk.StartAt) && now.Before(talk.EndAt) {
			return &talk, i
		}
	}
	return &Talk{}, 0
}

func (ts Talks) GetNextTalk(nextTalkListNum int) Talk {
	return ts[nextTalkListNum]
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
