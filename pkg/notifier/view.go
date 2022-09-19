package notifier

import (
	"github.com/slack-go/slack"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

func viewOnlineSession(talk model.Talk) slack.WebhookMessage {
	// TODO (#6)
	return slack.WebhookMessage{Text: "online session"}
}

func viewRecordingSession(talk model.Talk) slack.WebhookMessage {
	// TODO (#6)
	return slack.WebhookMessage{Text: "recording session"}
}

func viewCommercial(talk model.Talk) slack.WebhookMessage {
	// TODO (#6)
	return slack.WebhookMessage{Text: "commercial"}
}

func viewOpening(talk model.Talk) slack.WebhookMessage {
	// TODO (#6)
	return slack.WebhookMessage{Text: "opening"}
}

func viewEnding(talk model.Talk) slack.WebhookMessage {
	// TODO (#6)
	return slack.WebhookMessage{Text: "ending"}
}
