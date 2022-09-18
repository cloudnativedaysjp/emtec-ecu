package lib

import "context"

func (c *PrimitiveClient) ListTalks(ctx context.Context, eventAbbr string, trackId int32) (ListTalksResp, error) {
	// TODO
	return nil, nil
}

func (c *PrimitiveClient) UpdateTalks(ctx context.Context, talkId int32, onAir bool) error {
	// TODO
	return nil
}
