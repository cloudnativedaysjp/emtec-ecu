package lib

import "context"

type DreamkastApi interface {
	GenerateAuth0Token(ctx context.Context, auth0Domain, auth0ClientId, auth0ClientSecret, auth0Audience string) error
	ListTracks(ctx context.Context, eventAbbr string) (ListTracksResp, error)
	ListTalks(ctx context.Context, eventAbbr string, trackId int32) (ListTalksResp, error)
	UpdateTalks(ctx context.Context, talkId int32, onAir bool) error
}
