package notifier

import (
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/slack-go/slack"

	seaman_api "github.com/cloudnativedaysjp/seaman/seaman/api"

	"github.com/cloudnativedaysjp/emtec-ecu/pkg/model"
)

func talkBlocks(t model.Talk, titleMsg string, accessory *slack.Accessory) []interface{} {
	divider := map[string]interface{}{
		"type": "divider",
	}
	title := map[string]interface{}{
		"type": "context",
		"elements": []interface{}{
			map[string]interface{}{
				"emoji": true,
				"type":  "plain_text",
				"text":  titleMsg,
			},
		},
	}
	summary := map[string]interface{}{
		"type": "section",
		"fields": []interface{}{
			map[string]interface{}{
				"type":  "plain_text",
				"text":  fmt.Sprintf("Track %s", t.TrackName),
				"emoji": true,
			},
			map[string]interface{}{
				"type": "plain_text",
				"text": fmt.Sprintf("%s - %s",
					t.StartAt.Format("15:04"),
					t.EndAt.Format("15:04"),
				),
				"emoji": true,
			},
			map[string]interface{}{
				"type":  "plain_text",
				"text":  fmt.Sprintf("Type: %s", t.GetTalkTypeName()),
				"emoji": true,
			},
			map[string]interface{}{
				"emoji": true,
				"type":  "plain_text",
				"text": fmt.Sprintf("Speaker: %s",
					strings.Join(t.SpeakerNames, ", ")),
			},
		},
	}
	link := map[string]interface{}{
		"type": "section",
		"text": map[string]interface{}{
			"type": "mrkdwn",
			"text": fmt.Sprintf("Title: <%s/%s/talks/%d|%s>",
				eventUrlBase, t.EventAbbr, t.Id, t.TalkName),
		},
	}
	if accessory != nil {
		link["accessory"] = accessory
	}

	return []interface{}{
		divider,
		title,
		summary,
		link,
	}
}

func ViewNextSessionWillBegin(m *model.NotificationOnDkTimetable) slack.Msg {
	result, _ := viewNextSessionWillBegin(m)
	return result
}

func viewNextSessionWillBegin(m *model.NotificationOnDkTimetable) (slack.Msg, error) {
	currentTalk := m.Current()
	nextTalk := m.Next()

	// if currentTalk or nextTalk is on-demand session, create "switching" button
	accessory := &slack.Accessory{}
	if currentTalk.IsOnDemand() || nextTalk.IsOnDemand() {
		accessory = &slack.Accessory{
			ButtonElement: &slack.ButtonBlockElement{
				Type:     slack.METButton,
				ActionID: seaman_api.ActIdEmtec_SceneNext,
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

	blocks := []interface{}{
		map[string]interface{}{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": "*Next Scene will begin*",
			},
		},
	}
	// if currentTalk isn't empty, construct Block on CurrentTalk for Slack Msg
	if diff := cmp.Diff(currentTalk, model.Talk{}); diff != "" {
		blocks = append(blocks, talkBlocks(currentTalk, "Current Talk", nil)...)
	}
	// construct Block on NextTalk for Slack Msg
	blocks = append(blocks, talkBlocks(nextTalk, "Next Talk", accessory)...)

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
	currentTalk := m.Current()
	nextTalk := m.Next()

	blocks := []interface{}{
		map[string]interface{}{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": "*Scene was moved to next automatically*",
			},
		},
	}
	// if currentTalk isn't empty, construct Block on CurrentTalk for Slack Msg
	if diff := cmp.Diff(currentTalk, model.Talk{}); diff != "" {
		blocks = append(blocks, talkBlocks(currentTalk, "Previous Talk", nil)...)
	}
	// construct Block on NextTalk for Slack Msg
	blocks = append(blocks, talkBlocks(nextTalk, "Next Talk", nil)...)

	return castFromMapToMsg(
		map[string]interface{}{
			"blocks": blocks,
		},
	)
}
