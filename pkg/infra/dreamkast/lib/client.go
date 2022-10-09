package lib

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/utils"
)

type DreamkastClient interface {
	GenerateAuth0Token(ctx context.Context, auth0Domain, auth0ClientId, auth0ClientSecret, auth0Audience string) error
	ListTracks(ctx context.Context, eventAbbr string) (ListTracksResp, error)
	ListTalks(ctx context.Context, eventAbbr string, trackId int32) (ListTalksResp, error)
	UpdateTalks(ctx context.Context, talkId int32, onAir bool) error
}

type DreamkastClientImpl struct {
	client        *http.Client
	dkEndpointUrl url.URL
	auth0Token    string
}

func NewClient(dkEndpointUrl string) (DreamkastClient, error) {
	dkUrl, err := url.Parse(dkEndpointUrl)
	if err != nil {
		return nil, err
	}
	return &DreamkastClientImpl{client: http.DefaultClient, dkEndpointUrl: *dkUrl}, nil
}

func (c *DreamkastClientImpl) GenerateAuth0Token(ctx context.Context,
	auth0Domain, auth0ClientId, auth0ClientSecret, auth0Audience string,
) error {
	if c.auth0Token != "" {
		return nil
	}
	logger := utils.GetLogger(ctx)
	if auth0Domain == "" || auth0ClientId == "" || auth0ClientSecret == "" {
		logger.Info("auth0Domain or auth0ClientId or auth0ClientSecret was not set. " +
			"skipped to generate Auth0 Token")
		return nil
	}

	url, err := url.Parse(fmt.Sprintf("https://%s/oauth/token", auth0Domain))
	if err != nil {
		return err
	}
	payload := strings.NewReader("grant_type=client_credentials" +
		"&client_id=" + auth0ClientId +
		"&client_secret=" + auth0ClientSecret +
		"&audience=" + auth0Audience,
	)
	req, err := http.NewRequest("POST", url.String(), payload)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	c.auth0Token = string(body)
	return nil
}
