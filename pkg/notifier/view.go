package notifier

import (
	"github.com/slack-go/slack"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

func viewOnlineSession(talk model.Talk) slack.Msg {
	// TODO (#6)
	return slack.Msg{Text: "online session"}
}

func viewRecordingSession(talk model.Talk) slack.Msg {
	// TODO (#6)
	return slack.Msg{Text: "recording session"}
}

func viewCommercial(talk model.Talk) slack.Msg {
	// TODO (#6)
	return slack.Msg{Text: "commercial"}
}

func viewOpening(talk model.Talk) slack.Msg {
	// TODO (#6)
	return slack.Msg{Text: "opening"}
}

func viewEnding(talk model.Talk) slack.Msg {
	// TODO (#6)
	return slack.Msg{Text: "ending"}
}
