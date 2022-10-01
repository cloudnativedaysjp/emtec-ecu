package model

import (
	"fmt"
	"testing"
	"time"
)

const ISO8601ExtendedLayout = "2006-01-02T15:04:05"

func TestTalk_GetTalkType(t *testing.T) {
	talk := &Talk{}
	tests := []struct {
		name               string
		title              string
		presentationMethod interface{}
		want               TalkType
		wantErr            bool
	}{
		{
			name:               "opening",
			title:              "Opening",
			presentationMethod: nil,
			want:               3,
			wantErr:            false,
		},
		{
			name:               "commercial",
			title:              "休憩",
			presentationMethod: nil,
			want:               5,
			wantErr:            false,
		},
		{
			name:               "closing",
			title:              "Closing",
			presentationMethod: nil,
			want:               4,
			wantErr:            false,
		},
		{
			name:               "online session",
			title:              "CNDT",
			presentationMethod: "オンライン登壇",
			want:               1,
			wantErr:            false,
		},
		{
			name:               "recording session",
			title:              "recording",
			presentationMethod: "事前収録",
			want:               2,
			wantErr:            false,
		},
		{
			name:               "error",
			title:              "error",
			presentationMethod: "error",
			want:               0,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.presentationMethod == nil {
				got, err := talk.GetTalkType(tt.title, nil)
				if (err != nil) != tt.wantErr {
					t.Errorf("Talk.GetTalkType() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("Talk.GetTalkType() = %v, want %v", got, tt.want)
				}
			} else {
				str := fmt.Sprintf("%v", tt.presentationMethod)
				got, err := talk.GetTalkType(tt.title, &str)
				if (err != nil) != tt.wantErr {
					t.Errorf("Talk.GetTalkType() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("Talk.GetTalkType() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func setTestTime(v string) time.Time {
	pt, err := time.Parse(ISO8601ExtendedLayout, v)
	if err != nil {
		panic(err)
	}
	return pt
}

func TestTalk_WillStartNextTalkSince(t *testing.T) {
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
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.talks.WillStartNextTalkSince()
			if got != tt.want {
				t.Errorf("Talk.GetTalkType() want %v", tt.want)
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
			args: &Talk{
				Id:           999,
				TalkName:     "talk999",
				TrackId:      1,
				TrackName:    "track1",
				EventAbbr:    "test event",
				SpeakerNames: []string{"speakerA", "speaker"},
				Type:         1,
				StartAt:      setTestTime("2022-10-01T12:00:00"),
				EndAt:        setTestTime("2022-10-01T12:30:00"),
			},
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
			args: &Talk{
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
			args: &Talk{
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.talks.GetNextTalk(tt.args)
			if err == nil && tt.wantErr {
				t.Errorf("Talk.GetTalkType() wantErr %v", tt.wantErr)
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
			startAt, endAt, err := talk.GetActualStartAtAndEndAt(tt.conferenceDayDate, tt.startAt, tt.endAt)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Talk.GetActualStartAtAndEndAt() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				return
			}
			if err != nil || startAt.Equal(tt.startAt) || endAt.Equal(tt.endAt) {
				t.Errorf("Talk.GetActualStartAtAndEndAt() error = %v, wantStatAt = %v, wantEndAt = %v,", err, tt.wantStartAt, tt.wantEndAt)
				return
			}

		})
	}
}
