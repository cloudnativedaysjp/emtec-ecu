package notifier

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"

	seaman_api "github.com/cloudnativedaysjp/seaman/seaman/api"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

func ViewNextSessionWillBegin(m *model.NotificationOnDkTimetable) slack.Msg {
	result, _ := viewNextSessionWillBegin(m)
	return result
}

func viewNextSessionWillBegin(m *model.NotificationOnDkTimetable) (slack.Msg, error) {
	currentTalk := m.Current()
	nextTalk := m.Next()

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

	return castFromMapToMsg(
		map[string]interface{}{
			"blocks": []interface{}{
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
			},
		},
	)
}
