package notifier

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/cloudnativedaysjp/emtec-ecu/pkg/model"
)

func Test_viewNextSessionWillBegin(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		expectedStr := `
{
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "*Next Scene will begin*"
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "context",
			"elements": [
				{
					"type": "plain_text",
					"text": "Current Talk",
					"emoji": true
				}
			]
		},
		{
			"type": "section",
			"fields": [
				{
					"type": "plain_text",
					"text": "Track A",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "10:00 - 11:00",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "Type: オンライン登壇",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "Speaker: kanata",
					"emoji": true
				}
			]
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Title: <https://event.cloudnativedays.jp/cndt2101/talks/10001|ものすごい発表>"
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "context",
			"elements": [
				{
					"type": "plain_text",
					"text": "Next Talk",
					"emoji": true
				}
			]
		},
		{
			"type": "section",
			"fields": [
				{
					"type": "plain_text",
					"text": "Track A",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "11:00 - 12:30",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "Type: 事前収録",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "Speaker: hoge, fuga",
					"emoji": true
				}
			]
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Title: <https://event.cloudnativedays.jp/cndt2101/talks/10002|さらにものすごい発表>"
			},
			"accessory": {
				"type": "button",
				"action_id": "emtec_scenenext",
				"value": "1__A",
				"text": {
					"type": "plain_text",
					"text": "Switching"
				},
				"style": "primary",
				"confirm": {
					"title": {
						"type": "plain_text",
						"text": "Move to Next Scene"
					},
					"text": {
						"type": "plain_text",
						"text": "Are you sure?"
					},
					"confirm": {
						"type": "plain_text",
						"text": "OK"
					},
					"deny": {
						"type": "plain_text",
						"text": "Cancel"
					}
				}
			}
		}
	]
}
`
		expected, err := castFromStringToMsg(expectedStr)
		if err != nil {
			t.Fatal(err)
		}

		got, err := viewNextSessionWillBegin(model.NewNotificationOnDkTimetable(
			model.Talk{
				Id:           10001,
				TalkName:     "ものすごい発表",
				TrackId:      1,
				TrackName:    "A",
				StartAt:      time.Date(2022, 10, 1, 10, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
				EndAt:        time.Date(2022, 10, 1, 11, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
				Type:         model.TalkType_OnlineSession,
				SpeakerNames: []string{"kanata"},
				EventAbbr:    "cndt2101",
			},
			model.Talk{
				Id:           10002,
				TalkName:     "さらにものすごい発表",
				TrackId:      1,
				TrackName:    "A",
				StartAt:      time.Date(2022, 10, 1, 11, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
				EndAt:        time.Date(2022, 10, 1, 12, 30, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
				Type:         model.TalkType_RecordingSession,
				SpeakerNames: []string{"hoge", "fuga"},
				EventAbbr:    "cndt2101",
			},
		))
		if err != nil {
			t.Errorf("error = %v", err)
			return
		}
		if diff := cmp.Diff(expected, got); diff != "" {
			t.Errorf(diff)
		}
	})
}

func Test_viewScenemovedToNext(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		expectedStr := `
{
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "*Scene was moved to next automatically*"
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "context",
			"elements": [
				{
					"type": "plain_text",
					"text": "Previous Talk",
					"emoji": true
				}
			]
		},
		{
			"type": "section",
			"fields": [
				{
					"type": "plain_text",
					"text": "Track A",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "10:00 - 11:00",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "Type: オンライン登壇",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "Speaker: kanata",
					"emoji": true
				}
			]
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Title: <https://event.cloudnativedays.jp/cndt2101/talks/10001|ものすごい発表>"
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "context",
			"elements": [
				{
					"type": "plain_text",
					"text": "Next Talk",
					"emoji": true
				}
			]
		},
		{
			"type": "section",
			"fields": [
				{
					"type": "plain_text",
					"text": "Track A",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "11:00 - 12:30",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "Type: 事前収録",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "Speaker: hoge, fuga",
					"emoji": true
				}
			]
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Title: <https://event.cloudnativedays.jp/cndt2101/talks/10002|さらにものすごい発表>"
			}
		}
	]
}
`
		expected, err := castFromStringToMsg(expectedStr)
		if err != nil {
			t.Fatal(err)
		}

		got, err := viewSceneMovedToNext(model.NewNotificationSceneMovedToNext(
			model.Talk{
				Id:           10001,
				TalkName:     "ものすごい発表",
				TrackId:      1,
				TrackName:    "A",
				StartAt:      time.Date(2022, 10, 1, 10, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
				EndAt:        time.Date(2022, 10, 1, 11, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
				Type:         model.TalkType_OnlineSession,
				SpeakerNames: []string{"kanata"},
				EventAbbr:    "cndt2101",
			},
			model.Talk{
				Id:           10002,
				TalkName:     "さらにものすごい発表",
				TrackId:      1,
				TrackName:    "A",
				StartAt:      time.Date(2022, 10, 1, 11, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
				EndAt:        time.Date(2022, 10, 1, 12, 30, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
				Type:         model.TalkType_RecordingSession,
				SpeakerNames: []string{"hoge", "fuga"},
				EventAbbr:    "cndt2101",
			},
		))
		if err != nil {
			t.Errorf("error = %v", err)
			return
		}
		if diff := cmp.Diff(expected, got); diff != "" {
			t.Errorf(diff)
		}
	})
}
