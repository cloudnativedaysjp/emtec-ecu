package server

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-logr/logr"
	emptypb "github.com/golang/protobuf/ptypes/empty"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/obsws"
	pb "github.com/cloudnativedaysjp/cnd-operation-server/pkg/ws-proxy/scheme"
)

type Controller struct {
	pb.UnimplementedTrackServiceServer
	pb.UnimplementedSceneServiceServer

	Logger logr.Logger
	ObsWs  map[int32]obsws.ClientIface
}

func (c *Controller) GetTrack(ctx context.Context, in *pb.GetTrackRequest) (*pb.GetTrackResponse, error) {
	ws, ok := c.ObsWs[in.TrackId]
	if !ok {
		return nil, fmt.Errorf("no such trackId %d", in.TrackId)
	}

	return &pb.GetTrackResponse{Track: &pb.Track{
		TrackId: in.TrackId,
		ObsHost: ws.GetHost(),
	}}, nil
}

func (c *Controller) ListTrack(ctx context.Context, in *emptypb.Empty) (*pb.ListTrackResponse, error) {
	var tracks []*pb.Track
	for trackId, ws := range c.ObsWs {
		tracks = append(tracks, &pb.Track{
			TrackId: trackId,
			ObsHost: ws.GetHost(),
		})
	}
	sort.SliceStable(tracks, func(i, j int) bool { return tracks[i].TrackId < tracks[j].TrackId })

	return &pb.ListTrackResponse{Track: tracks}, nil
}

func (c *Controller) ListScene(ctx context.Context, in *pb.ListSceneRequest) (*pb.ListSceneResponse, error) {
	ctx = logr.NewContext(ctx, c.Logger)
	ws, ok := c.ObsWs[in.TrackId]
	if !ok {
		return nil, fmt.Errorf("no such trackId %d", in.TrackId)
	}

	scenes, err := ws.ListScenes(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListSceneResponse{}
	for _, scene := range scenes {
		resp.Scene = append(resp.Scene, &pb.Scene{
			Name:             scene.Name,
			SceneIndex:       int32(scene.SceneIndex),
			IsCurrentProgram: scene.IsCurrentProgram,
		})
	}
	return resp, nil
}

func (c *Controller) MoveSceneToNext(ctx context.Context, in *pb.MoveSceneToNextRequest) (*emptypb.Empty, error) {
	ctx = logr.NewContext(ctx, c.Logger)
	ws, ok := c.ObsWs[in.TrackId]
	if !ok {
		return nil, fmt.Errorf("no such trackId %d", in.TrackId)
	}

	if err := ws.MoveSceneToNext(ctx); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
