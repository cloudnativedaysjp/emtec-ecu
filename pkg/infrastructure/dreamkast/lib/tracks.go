package lib

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
)

func (c *PrimitiveClient) ListTracks(ctx context.Context, eventAbbr string) (ListTracksResp, error) {
	url := c.dkEndpointUrl
	url.Path = filepath.Join(url.Path, "/v1/tracks")
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("eventAbbr", eventAbbr)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ListTracksResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
