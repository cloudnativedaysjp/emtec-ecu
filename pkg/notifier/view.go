package notifier

import (
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/slack-go/slack"

	seaman_api "github.com/cloudnativedaysjp/seaman/seaman/api"

	"github.com/cloudnativedaysjp/emtec-ecu/pkg/model"
)

func ViewNextSessionWillBegin(m *model.NotificationOnDkTimetable) slack.Msg {
	result, _ := viewNextSessionWillBegin(m)
	return result
}

func viewNextSessionWillBegin(m *model.NotificationOnDkTimetable) (slack.Msg, error) {
	currentTalk := m.Current()
	nextTalk := m.Next()
	var BlocksForCurrentTalk, BlocksForNextTalk []interface{}

	// if currentTalk or nextTalk is on-demand session, create "switching" button
	accessory := &slack.Accessory{}
	if currentTalk.IsOnDemand() || nextTalk.IsOnDemand() {
		accessory = &slack.Accessory{
			ButtonElement: &slack.ButtonBlockElement{
				Type:     slack.METButton,
				ActionID: seaman_api.ActIdBroadcast_SceneNext,
				Value:    seaman_api.Track{Id: m.TrackId(), Name: m.TrackName()}.String(),
				Text: &slack.TextBlockObject{
					Type: "plain_text",
					Text: "Switching",
				},
				Style: slack.StylePrimary,
				Confirm: &slack.ConfirmationBlockObject{
					Title: &slack.TextBlockObject{
						Type: "plain_text",
						Text: "Move to Next Scene",
					},
					Text: &slack.TextBlockObject{
						Type: "plain_text",
						Text: "Are you sure?",
					},
					Confirm: &slack.TextBlockObject{
						Type: "plain_text",
						Text: "OK",
					},
					Deny: &slack.TextBlockObject{
						Type: "plain_text",
						Text: "Cancel",
					},
				},
			},
		}
	}

	// if currentTalk isn't empty, construct Block on CurrentTalk for Slack Msg
	if diff := cmp.Diff(currentTalk, model.Talk{}); diff != "" {
		BlocksForCurrentTalk = []interface{}{
			map[string]interface{}{
				"type": "divider",
			},
			map[string]interface{}{
				"type": "context",
				"elements": []interface{}{
					map[string]interface{}{
						"emoji": true,
						"type":  "plain_text",
						"text":  "Current Talk",
					},
				},
			},
			map[string]interface{}{
				"type": "section",
				"fields": []interface{}{
					map[string]interface{}{
						"type":  "plain_text",
						"text":  fmt.Sprintf("Track %s", currentTalk.TrackName),
						"emoji": true,
					},
					map[string]interface{}{
						"type": "plain_text",
						"text": fmt.Sprintf("%s - %s",
							currentTalk.StartAt.Format("15:04"),
							currentTalk.EndAt.Format("15:04"),
						),
						"emoji": true,
					},
					map[string]interface{}{
						"type":  "plain_text",
						"text":  fmt.Sprintf("Type: %s", currentTalk.GetTalkTypeName()),
						"emoji": true,
					},
					map[string]interface{}{
						"emoji": true,
						"type":  "plain_text",
						"text": fmt.Sprintf("Speaker: %s",
							strings.Join(currentTalk.SpeakerNames, ", ")),
					},
				},
			},
			map[string]interface{}{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("Title: <%s/%s/talks/%d|%s>",
						eventUrlBase, currentTalk.EventAbbr, currentTalk.Id, currentTalk.TalkName),
				},
			},
		}
	}

	// construct Block on NextTalk for Slack Msg
	BlocksForNextTalk = []interface{}{
		map[string]interface{}{
			"type": "divider",
		},
		map[string]interface{}{
			"type": "context",
			"elements": []interface{}{
				map[string]interface{}{
					"type":  "plain_text",
					"text":  "Next Talk",
					"emoji": true,
				},
			},
		},
		map[string]interface{}{
			"type": "section",
			"fields": []interface{}{
				map[string]interface{}{
					"type":  "plain_text",
					"text":  fmt.Sprintf("Track %s", nextTalk.TrackName),
					"emoji": true,
				},
				map[string]interface{}{
					"type": "plain_text",
					"text": fmt.Sprintf("%s - %s",
						nextTalk.StartAt.Format("15:04"),
						nextTalk.EndAt.Format("15:04"),
					),
					"emoji": true,
				},
				map[string]interface{}{
					"type":  "plain_text",
					"text":  fmt.Sprintf("Type: %s", nextTalk.GetTalkTypeName()),
					"emoji": true,
				},
				map[string]interface{}{
					"type": "plain_text",
					"text": fmt.Sprintf("Speaker: %s",
						strings.Join(nextTalk.SpeakerNames, ", ")),
					"emoji": true,
				},
			},
		},
		map[string]interface{}{
			"accessory": accessory,
			"type":      "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": fmt.Sprintf("Title: <%s/%s/talks/%d|%s>",
					eventUrlBase, nextTalk.EventAbbr, nextTalk.Id, nextTalk.TalkName),
			},
		},
	}

	blocks := []interface{}{
		map[string]interface{}{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": "*Next Scene will begin*",
			},
		},
	}
	blocks = append(blocks, BlocksForCurrentTalk...)
	blocks = append(blocks, BlocksForNextTalk...)

	return castFromMapToMsg(
		map[string]interface{}{
			"blocks": blocks,
		},
	)
}

func ViewSceneMovedToNext(m *model.NotificationSceneMovedToNext) slack.Msg {
	result, _ := viewSceneMovedToNext(m)
	return result
}

func viewSceneMovedToNext(m *model.NotificationSceneMovedToNext) (slack.Msg, error) {
	nextTalk := m.Next()
	return castFromMapToMsg(
		map[string]interface{}{
			"blocks": []interface{}{
				map[string]interface{}{
					"type": "section",
					"text": map[string]interface{}{
						"type": "mrkdwn",
						"text": "*Scene was moved to next automatically*",
					},
				},
				map[string]interface{}{
					"type": "divider",
				},
				map[string]interface{}{
					"type": "context",
					"elements": []interface{}{
						map[string]interface{}{
							"emoji": true,
							"type":  "plain_text",
							"text":  "Current Talk",
						},
					},
				},
				map[string]interface{}{
					"type": "section",
					"fields": []interface{}{
						map[string]interface{}{
							"type":  "plain_text",
							"text":  fmt.Sprintf("Track %s", nextTalk.TrackName),
							"emoji": true,
						},
						map[string]interface{}{
							"type": "plain_text",
							"text": fmt.Sprintf("%s - %s",
								nextTalk.StartAt.Format("15:04"),
								nextTalk.EndAt.Format("15:04"),
							),
							"emoji": true,
						},
						map[string]interface{}{
							"type":  "plain_text",
							"text":  fmt.Sprintf("Type: %s", nextTalk.GetTalkTypeName()),
							"emoji": true,
						},
						map[string]interface{}{
							"emoji": true,
							"type":  "plain_text",
							"text": fmt.Sprintf("Speaker: %s",
								strings.Join(nextTalk.SpeakerNames, ", ")),
						},
					},
				},
				map[string]interface{}{
					"type": "section",
					"text": map[string]interface{}{
						"type": "mrkdwn",
						"text": fmt.Sprintf("Title: <%s/%s/talks/%d|%s>",
							eventUrlBase, nextTalk.EventAbbr, nextTalk.Id, nextTalk.TalkName),
					},
				},
			},
		},
	)
}
