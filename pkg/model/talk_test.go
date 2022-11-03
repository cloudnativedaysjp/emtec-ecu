package model

import (
	"testing"
	"time"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/testutils"
)

const (
	test_ISO8601ExtendedLayout     = "2006-01-02T15:04:05"
	test_howManyMinutesUntilNotify = 5 * time.Minute
)

func TestTalk_GetTalkType(t *testing.T) {
	talk := &Talk{}
	tests := []struct {
		name               string
		title              string
		presentationMethod *string
		want               TalkType
		wantErr            bool
	}{
		{
			name:               "opening",
			title:              "Opening",
			presentationMethod: nil,
			want:               TalkType_Opening,
			wantErr:            false,
		},
		{
			name:               "commercial",
			title:              "休憩",
			presentationMethod: nil,
			want:               TalkType_Commercial,
			wantErr:            false,
		},
		{
			name:               "closing",
			title:              "Closing",
			presentationMethod: nil,
			want:               TalkType_Ending,
			wantErr:            false,
		},
		{
			name:               "online session",
			title:              "CNDT",
			presentationMethod: testutils.ToPointer(t, "オンライン登壇"),
			want:               TalkType_OnlineSession,
			wantErr:            false,
		},
		{
			name:               "recording session",
			title:              "recording",
			presentationMethod: testutils.ToPointer(t, "事前収録"),
			want:               TalkType_RecordingSession,
			wantErr:            false,
		},
		{
			name:               "error",
			title:              "error",
			presentationMethod: testutils.ToPointer(t, "error"),
			want:               0,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := talk.GetTalkType(tt.title, tt.presentationMethod)
			if (err != nil) != tt.wantErr {
				t.Errorf("Talk.GetTalkType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Talk.GetTalkType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setTestTime(v string) time.Time {
	pt, err := time.Parse(test_ISO8601ExtendedLayout, v)
	if err != nil {
		panic(err)
	}
	return pt
}

func TestTalk_IsStartNextTalkSoon(t *testing.T) {
	nowFunc = func() time.Time {
		return time.Date(2022, 10, 01, 12, 27, 00, 0, time.UTC)
	}
	tests := []struct {
		name  string
		talks Talks
		want  bool
	}{
		{
			name:  "no talk",
			talks: Talks{},
			want:  false,
		},
		{
			name: "not found next talk",
			talks: Talks{
				Talk{
					Id:           1,
					TalkName:     "talk1",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T11:00:00"),
					EndAt:        setTestTime("2022-10-01T12:00:00"),
				},
				Talk{
					Id:           2,
					TalkName:     "talk2",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T12:00:00"),
					EndAt:        setTestTime("2022-10-01T13:00:00"),
				},
				Talk{
					Id:           3,
					TalkName:     "talk3",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T13:00:00"),
					EndAt:        setTestTime("2022-10-01T14:00:00"),
				},
			},
			want: false,
		},
		{
			name: "found next talk",
			talks: Talks{
				Talk{
					Id:           1,
					TalkName:     "talk1",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T11:00:00"),
					EndAt:        setTestTime("2022-10-01T11:30:00"),
				},
				Talk{
					Id:           2,
					TalkName:     "talk2",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T11:30:00"),
					EndAt:        setTestTime("2022-10-01T12:00:00"),
				},
				Talk{
					Id:           3,
					TalkName:     "talk3",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T12:00:00"),
					EndAt:        setTestTime("2022-10-01T12:30:00"),
				},
				Talk{
					Id:           4,
					TalkName:     "talk4",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T12:30:00"),
					EndAt:        setTestTime("2022-10-01T13:00:00"),
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.talks.IsStartNextTalkSoon(test_howManyMinutesUntilNotify)
			if got != tt.want {
				t.Errorf("Talk.HasNotify() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestTalk_GetCurrentTalk(t *testing.T) {
	nowFunc = func() time.Time {
		return time.Date(2022, 10, 01, 12, 27, 00, 0, time.UTC)
	}
	tests := []struct {
		name    string
		talks   Talks
		wantId  int
		wantErr bool
	}{
		{
			name:    "no talk",
			talks:   Talks{},
			wantErr: true,
		},
		{
			name: "not found current talk",
			talks: Talks{
				Talk{
					Id:           1,
					TalkName:     "talk1",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T01:00:00"),
					EndAt:        setTestTime("2022-10-01T02:00:00"),
				},
				Talk{
					Id:           2,
					TalkName:     "talk2",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T02:00:00"),
					EndAt:        setTestTime("2022-10-01T03:00:00"),
				},
				Talk{
					Id:           3,
					TalkName:     "talk3",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T03:00:00"),
					EndAt:        setTestTime("2022-10-01T04:00:00"),
				},
			},
			wantErr: true,
		},
		{
			name: "found current talk",
			talks: Talks{
				Talk{
					Id:           1,
					TalkName:     "talk1",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T11:00:00"),
					EndAt:        setTestTime("2022-10-01T11:30:00"),
				},
				Talk{
					Id:           2,
					TalkName:     "talk2",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T11:30:00"),
					EndAt:        setTestTime("2022-10-01T12:00:00"),
				},
				Talk{
					Id:           3,
					TalkName:     "talk3",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T12:00:00"),
					EndAt:        setTestTime("2022-10-01T12:30:00"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.talks.GetCurrentTalk()
			if err == nil && tt.wantErr {
				t.Errorf("Talk.GetTalkType() wantErr %v", tt.wantErr)
			}
		})
	}
}

func TestTalk_GetNextTalk(t *testing.T) {
	nowFunc = func() time.Time {
		return time.Date(2022, 10, 01, 12, 27, 00, 0, time.UTC)
	}
	tests := []struct {
		name    string
		args    *Talk
		talks   Talks
		wantId  int32
		wantErr bool
	}{
		{
			name:    "no talk",
			args:    nil,
			talks:   Talks{},
			wantErr: true,
		},
		{
			name: "not found next talk",
			talks: Talks{
				Talk{
					Id:           1,
					TalkName:     "talk1",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T01:00:00"),
					EndAt:        setTestTime("2022-10-01T02:00:00"),
				},
				Talk{
					Id:           2,
					TalkName:     "talk2",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T02:00:00"),
					EndAt:        setTestTime("2022-10-01T03:00:00"),
				},
				Talk{
					Id:           3,
					TalkName:     "talk3",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T03:00:00"),
					EndAt:        setTestTime("2022-10-01T04:00:00"),
				},
			},
			wantErr: true,
		},
		{
			name: "last talk",
			talks: Talks{
				Talk{
					Id:           1,
					TalkName:     "talk1",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T11:00:00"),
					EndAt:        setTestTime("2022-10-01T11:30:00"),
				},
				Talk{
					Id:           2,
					TalkName:     "talk2",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T11:30:00"),
					EndAt:        setTestTime("2022-10-01T12:00:00"),
				},
				Talk{
					Id:           3,
					TalkName:     "talk3",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T12:00:00"),
					EndAt:        setTestTime("2022-10-01T12:30:00"),
				},
			},
			wantErr: true,
		},
		{
			name: "found next talk",
			talks: Talks{
				Talk{
					Id:           1,
					TalkName:     "talk1",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T11:00:00"),
					EndAt:        setTestTime("2022-10-01T11:30:00"),
				},
				Talk{
					Id:           2,
					TalkName:     "talk2",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T11:30:00"),
					EndAt:        setTestTime("2022-10-01T12:00:00"),
				},
				Talk{
					Id:           3,
					TalkName:     "talk3",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T12:00:00"),
					EndAt:        setTestTime("2022-10-01T12:30:00"),
				},
				Talk{
					Id:           4,
					TalkName:     "talk4",
					TrackId:      1,
					TrackName:    "track1",
					EventAbbr:    "test event",
					SpeakerNames: []string{"speakerA", "speaker"},
					Type:         1,
					StartAt:      setTestTime("2022-10-01T12:30:00"),
					EndAt:        setTestTime("2022-10-01T13:00:00"),
				},
			},
			wantErr: false,
			wantId:  4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.talks.GetNextTalk()
			if (err == nil && tt.wantErr) || (err != nil && !tt.wantErr) {
				t.Errorf("Talk.GetTalkType() wantErr %v", tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.Id != int32(tt.wantId) {
				t.Errorf("Talk.GetTalkType() wantId %v", tt.wantId)
				return
			}
		})
	}
}

func TestTalk_GetActualStartAtAndEndAt(t *testing.T) {
	talk := &Talk{}
	tests := []struct {
		name              string
		conferenceDayDate string
		startAt           time.Time
		endAt             time.Time
		wantStartAt       time.Time
		wantEndAt         time.Time
		wantErr           bool
	}{
		{
			name:              "invalid conference date format",
			conferenceDayDate: "2022/10/01",
			startAt:           setTestTime("2000-01-01T12:00:00"),
			endAt:             setTestTime("2022-01-01T12:00:00"),
			wantStartAt:       setTestTime("2022-10-01T12:00:00"),
			wantEndAt:         setTestTime("2022-10-01T12:00:00"),
			wantErr:           true,
		},
		{
			name:              "ok",
			conferenceDayDate: "2022-10-01",
			startAt:           setTestTime("2000-01-01T12:00:00"),
			endAt:             setTestTime("2022-01-01T13:00:00"),
			wantStartAt:       setTestTime("2022-10-01T12:00:00"),
			wantEndAt:         setTestTime("2022-10-01T13:00:00"),
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStartAt, gotEndAt, err := talk.GetActualStartAtAndEndAt(tt.conferenceDayDate, tt.startAt, tt.endAt)
			if (err == nil && tt.wantErr) || (err != nil && !tt.wantErr) {
				t.Errorf("Talk.GetActualStartAtAndEndAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !gotStartAt.Equal(tt.wantStartAt) || !gotEndAt.Equal(tt.wantEndAt) {
				t.Errorf("Talk.GetActualStartAtAndEndAt() error = %v, wantStatAt = %v, wantEndAt = %v,", err, tt.wantStartAt, tt.wantEndAt)
				return
			}

		})
	}
}
