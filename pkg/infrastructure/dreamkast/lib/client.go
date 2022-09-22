package lib

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type PrimitiveClient struct {
	client        *http.Client
	dkEndpointUrl url.URL
	auth0Token    string
}

func NewClient(dkEndpointUrl string) (DreamkastApi, error) {
	dkUrl, err := url.Parse(dkEndpointUrl)
	if err != nil {
		return nil, err
	}
	return &PrimitiveClient{client: http.DefaultClient, dkEndpointUrl: *dkUrl}, nil
}

func (c *PrimitiveClient) GenerateAuth0Token(
	auth0Domain, auth0ClientId, auth0ClientSecret, auth0Audience string,
) error {
	if c.auth0Token != "" {
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
