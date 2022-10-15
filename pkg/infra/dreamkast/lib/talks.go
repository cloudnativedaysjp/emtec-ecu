package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
)

func (c *DreamkastClientImpl) ListTalks(ctx context.Context, eventAbbr string, trackId int32) (ListTalksResp, error) {
	url := c.dkEndpointUrl
	url.Path = filepath.Join(url.Path, "/v1/talks")
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("eventAbbr", eventAbbr)
	q.Add("trackId", strconv.Itoa(int(trackId)))
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ListTalksResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// TODO: set Token to header
// TODO: write test
func (c *DreamkastClientImpl) UpdateTalk(ctx context.Context, talkId int32, onAir bool) error {
	url := c.dkEndpointUrl
	url.Path = filepath.Join(url.Path, "/v1/talks", strconv.Itoa(int(talkId)))
	reqBody, err := json.Marshal(&UpdateTalksReq{OnAir: onAir})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	if _, err := c.client.Do(req); err != nil {
		return err
	}
	return nil
}
