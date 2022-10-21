package model

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/infra"
)

type TalkType int32

const (
	TalkType_OnlineSession TalkType = iota + 1
	TalkType_RecordingSession
	TalkType_Opening
	TalkType_Ending
	TalkType_Commercial

	dateLayout = "2006-01-02"
)

//
// Talks
//

type Talks []Talk

func (ts Talks) IsStartNextTalkSoon(untilNotify time.Duration) bool {
	now := nowFunc()
	nextTalk, err := ts.GetNextTalk()
	if err != nil {
		return false
	}
	return nextTalk.StartAt.Sub(now) <= untilNotify
}

func (ts Talks) HasNotify(ctx context.Context, rc *infra.RedisClient, untilNotify time.Duration) (bool, error) {
	result := rc.Client.Get(ctx, infra.NextTalkNotificationKey)
	if result.Err() != nil {
		return false, result.Err()
	}
	if result != nil {
		return false, nil
	}
	return true, nil
}

func (ts Talks) GetCurrentTalk() (*Talk, error) {
	if len(ts) == 0 {
		return nil, fmt.Errorf("Talks is empty")
	}
	now := nowFunc()
	for _, talk := range ts {
		if now.After(talk.StartAt) && now.Before(talk.EndAt) {
			return &talk, nil
		}
	}
	return nil, fmt.Errorf("current talk not found")
}

func (ts Talks) GetNextTalk() (*Talk, error) {
	if len(ts) == 0 {
		return nil, fmt.Errorf("talks is empty")
	}

	currentTalk, err := ts.GetCurrentTalk()
	if err != nil {
		// When currentTalk is none,
		// if now is before event then return Talks[0] else raise an error
		now := nowFunc()
		lastTalk := ts[len(ts)-1]
		if now.After(lastTalk.EndAt) {
			return nil, fmt.Errorf("talks has already finished")
		}
		return &ts[0], nil
	}
	for i, talk := range ts {
		if talk.Id == currentTalk.Id {
			if i+1 == len(ts) {
				return nil, fmt.Errorf("current talk is last")
			}
			return &ts[i+1], nil
		}
	}
	return nil, fmt.Errorf("something Wrong")
}

//
// Talk
//

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

func (t Talk) GetTalkTypeName() string {
	var typeName string
	switch t.Type {
	case TalkType_OnlineSession:
		typeName = "オンライン登壇"
	case TalkType_RecordingSession:
		typeName = "事前収録"
	case TalkType_Opening:
		typeName = "Opening"
	case TalkType_Ending:
		typeName = "Closing"
	case TalkType_Commercial:
		typeName = "CM"
	}
	return typeName
}

func (t Talk) IsOnDemand() bool {
	return t.Type == TalkType_OnlineSession || t.Type == TalkType_Opening || t.Type == TalkType_Ending
}

func (t Talk) GetActualStartAtAndEndAt(conferenceDayDate string, startAt, endAt time.Time) (time.Time, time.Time, error) {
	cDate, err := time.Parse(dateLayout, conferenceDayDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return cDate.Add(time.Duration(startAt.Hour()*int(time.Hour) + startAt.Minute()*int(time.Minute) + startAt.Second()*int(time.Second))),
		cDate.Add(time.Duration(endAt.Hour()*int(time.Hour) + endAt.Minute()*int(time.Minute) + endAt.Second()*int(time.Second))),
		nil
}

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

//
// CurrentAndNextTalk
//

type CurrentAndNextTalk struct {
	Current Talk
	Next    Talk
}

func (m CurrentAndNextTalk) TrackId() int32 {
	if m.Current.TrackId != 0 {
		return m.Current.TrackId
	}
	return m.Next.TrackId
}

func (m CurrentAndNextTalk) TrackName() string {
	if m.Current.TrackName != "" {
		return m.Current.TrackName
	}
	return m.Next.TrackName
}
