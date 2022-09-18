package lib

import "time"

type ListTracksResp []struct {
	ID            int32     `json:"id"`
	Name          string    `json:"name"`
	VideoPlatform string    `json:"video_platform"`
	VideoID       *string   `json:"video_id"`
	ChannelArn    *string   `json:"channel_arn"`
	OnAirTalk     OnAirTalk `json:"on_air_talk"`
}

type OnAirTalk struct {
	ID            int32     `json:"id"`
	TalkID        int32     `json:"talk_id"`
	Site          *string   `json:"site"`
	Url           *string   `json:"url"`
	OnAir         bool      `json:"on_air"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	VideoID       string    `json:"video_id"`
	SlidoID       *string   `json:"slido_id"`
	VideoFileData *string   `json:"video_file_data"`
}

type ListTalksResp []struct {
	GetTalksResp `json:""`
	ConferenceID int32 `json:"conference_id"`
}

type GetTalksResp struct {
	ID                int32       `json:"id"`
	TrackID           int32       `json:"trackID"`
	VideoPlatform     interface{} `json:"videoPlatform"`
	VideoID           string      `json:"videoID"`
	Title             string      `json:"title"`
	Abstract          string      `json:"abstract"`
	Speakers          []Speaker   `json:"speakers"`
	DayID             int32       `json:"dayID"`
	ShowOnTimetable   bool        `json:"showOnTimetable"`
	StartTime         time.Time   `json:"startTime"`
	EndTime           time.Time   `json:"endTime"`
	TalkDuration      int32       `json:"talkDuration"`
	TalkDifficulty    string      `json:"talkDifficulty"`
	TalkCategory      string      `json:"talkCategory"`
	OnAir             bool        `json:"onAir"`
	DocumentURL       string      `json:"documentUrl"`
	ConferenceDayID   int32       `json:"conferenceDayID"`
	ConferenceDayDate string      `json:"conferenceDayDate"`
	StartOffset       int32       `json:"startOffset"`
	EndOffset         int32       `json:"endOffset"`
	ActualStartTime   time.Time   `json:"actualStartTime"`
	ActualEndTime     time.Time   `json:"actualEndTime"`
}

type Speaker struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

//
// Error Types
//

type ErrorUnauthorized struct{}

func (e ErrorUnauthorized) Error() string {
	return "Unauthorized"
}
