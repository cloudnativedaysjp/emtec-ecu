package model

import (
	"time"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/utils"
)

type TalkType int32

const (
	TalkType_OnlineSession TalkType = iota + 1
	TalkType_RecordingSession
	TalkType_Opening
	TalkType_Ending
	TalkType_Commercial
)

type Talks []Talk

func (t Talk) ConvertTalkType(title string, presentationMethod *string) TalkType {
	switch {
	case presentationMethod == nil:
		switch title {
		case "Opening":
			return TalkType_Opening
		case "休憩":
			return TalkType_Commercial
		case "Closing":
			return TalkType_Ending
		}
	case *presentationMethod == "オンライン登壇":
		return TalkType_OnlineSession
	case *presentationMethod == "事前収録":
		return TalkType_RecordingSession
	}
	return 0
}

func (t Talk) GetTalkType(title string, presentationMethod *string) TalkType {
	return t.ConvertTalkType(title, presentationMethod)
}

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

func (ts Talks) GetNextTalk(currentTalkListNum int) Talk {
	return ts[currentTalkListNum+1]
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
