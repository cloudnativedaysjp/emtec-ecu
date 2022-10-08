package notifier

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"

	seaman_api "github.com/cloudnativedaysjp/seaman/seaman/api"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

func ViewSession(m model.CurrentAndNextTalk) slack.Msg {
	result, _ := viewSession(m)
	return result
}

func viewSession(m model.CurrentAndNextTalk) (slack.Msg, error) {
	currentTalk := m.Current
	nextTalk := m.Next

	accessory := make(map[string]interface{})
	if currentTalk.IsOnDemand() || nextTalk.IsOnDemand() {
		accessory = map[string]interface{}{
			"action_id": "multi_static_select-action",
			"type":      "multi_static_select",
			"placeholder": map[string]interface{}{
				"text":  "switching",
				"emoji": true,
				"type":  "plain_text",
			},
			"options": []interface{}{
				map[string]interface{}{
					"text": map[string]interface{}{
						"type":  "plain_text",
						"text":  "シーンを切り替える",
						"emoji": true,
					},
					"value": seaman_api.ActIdBroadcast_SceneNext,
				},
			},
		}
	}

	return castFromMapToMsg(
		map[string]interface{}{
			"blocks": []interface{}{
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
						"text": fmt.Sprintf("<%s/%s/talks/%d|Title: %s>",
							eventUrlBase, currentTalk.EventAbbr, currentTalk.Id, currentTalk.TalkName),
					},
				},
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
						"text": fmt.Sprintf("<%s/%s/talks/%d|Title: %s>",
							eventUrlBase, nextTalk.EventAbbr, nextTalk.Id, nextTalk.TalkName),
					},
				},
				map[string]interface{}{
					"type": "divider",
				},
			},
		},
	)
}
