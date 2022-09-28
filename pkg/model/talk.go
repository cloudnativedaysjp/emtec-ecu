package model

import (
	"fmt"
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

func (t Talk) convertTalkType(title string, presentationMethod *string) (TalkType, error) {
	switch {
	case presentationMethod == nil:
		switch title {
		case "Opening":
			return TalkType_Opening, nil
		case "休憩":
			return TalkType_Commercial, nil
		case "Closing":
			return TalkType_Ending, nil
		}
	case *presentationMethod == "オンライン登壇":
		return TalkType_OnlineSession, nil
	case *presentationMethod == "事前収録":
		return TalkType_RecordingSession, nil
	}
	return 0, fmt.Errorf("model.convertTalkType not found. title: %s, presentationMethod: %s", title, *presentationMethod)
}

func (t Talk) GetTalkType(title string, presentationMethod *string) (TalkType, error) {
	return t.convertTalkType(title, presentationMethod)
}

func (ts Talks) WillStartNextTalkSince() bool {
	now := nowFunc()
	for _, talk := range ts {
		if now.After(talk.StartAt) {
			diffTime := time.Duration(talk.EndAt.Sub(now).Minutes())
			if 0 < diffTime && diffTime <= utils.HowManyMinutesUntilNotify {
				return true
			}
		}
	}
	return false
}

func (ts Talks) GetCurrentTalk() (*Talk, error) {
	now := nowFunc()
	for _, talk := range ts {
		if now.After(talk.StartAt) && now.Before(talk.EndAt) {
			return &talk, nil
		}
	}
	return nil, fmt.Errorf("Current talk not found")
}

func (ts Talks) GetNextTalk(currentTalk *Talk) (*Talk, error) {
	for i, talk := range ts {
		if talk.Id == currentTalk.Id {
			if i == len(ts) {
				fmt.Errorf("This talk is last")
			}
			return &ts[i+1], nil
		}
	}
	return nil, fmt.Errorf("Next talk not found")
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
