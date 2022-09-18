package cmd

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/cloudnativedaysjp/cnd-operation-server/pkg/ws-proxy/scheme"
)

var (
	sceneClient pb.SceneServiceClient
	trackClient pb.TrackServiceClient
)

func createClient(addr string) error {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // TODO (#7)
	)
	if err != nil {
		return err
	}

	sceneClient = pb.NewSceneServiceClient(conn)
	trackClient = pb.NewTrackServiceClient(conn)
	return nil
}
