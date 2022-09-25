package server

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-logr/logr"
	emptypb "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/obsws"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/sharedmem"
	pb "github.com/cloudnativedaysjp/cnd-operation-server/pkg/ws-proxy/schema"
)

type Controller struct {
	pb.UnimplementedTrackServiceServer
	pb.UnimplementedSceneServiceServer
	pb.UnimplementedDebugServiceServer

	Logger      logr.Logger
	ObsWsMap    map[int32]obsws.ClientIface
	MemWriter   sharedmem.WriterIface
	MemDebugger sharedmem.DebuggerIface
}

/* TrackService */

func (c *Controller) GetTrack(ctx context.Context, in *pb.GetTrackRequest) (*pb.GetTrackResponse, error) {
	ws, ok := c.ObsWsMap[in.TrackId]
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
	for trackId, ws := range c.ObsWsMap {
		tracks = append(tracks, &pb.Track{
			TrackId: trackId,
			ObsHost: ws.GetHost(),
		})
	}
	sort.SliceStable(tracks, func(i, j int) bool { return tracks[i].TrackId < tracks[j].TrackId })

	return &pb.ListTrackResponse{Track: tracks}, nil
}

func (c *Controller) EnableAutomation(ctx context.Context, in *pb.SwitchAutomationRequest) (*emptypb.Empty, error) {
	if err := c.MemWriter.SetDisableAutomation(in.TrackId, false); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (c *Controller) DisableAutomation(ctx context.Context, in *pb.SwitchAutomationRequest) (*emptypb.Empty, error) {
	if err := c.MemWriter.SetDisableAutomation(in.TrackId, true); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

/* SceneService */

func (c *Controller) ListScene(ctx context.Context, in *pb.ListSceneRequest) (*pb.ListSceneResponse, error) {
	ctx = logr.NewContext(ctx, c.Logger)
	ws, ok := c.ObsWsMap[in.TrackId]
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
	ws, ok := c.ObsWsMap[in.TrackId]
	if !ok {
		return nil, fmt.Errorf("no such trackId %d", in.TrackId)
	}

	if err := ws.MoveSceneToNext(ctx); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

/* DebugService */

func (c *Controller) ListSharedmem(context.Context, *emptypb.Empty) (*pb.ListSharedmemResponse, error) {
	disabledMap := c.MemDebugger.ListAutomation()
	talksMap := make(map[int32]*pb.TalksModel)
	for k, v := range c.MemDebugger.ListTalks() {
		var talks []*pb.TalkModel
		for _, talk := range v {
			talks = append(talks, &pb.TalkModel{
				Id:           talk.Id,
				TalkName:     talk.TalkName,
				TrackId:      talk.TrackId,
				TrackName:    talk.TrackName,
				EventAbbr:    talk.EventAbbr,
				SpeakerNames: talk.SpeakerNames,
				Type:         int32(talk.Type),
				StartAt:      timestamppb.New(talk.StartAt),
				EndAt:        timestamppb.New(talk.EndAt),
			})
		}
		talksMap[k] = &pb.TalksModel{Talks: talks}
	}
	return &pb.ListSharedmemResponse{
		TalksMap:    talksMap,
		DisabledMap: disabledMap,
	}, nil
}
