package cmd

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/cloudnativedaysjp/cnd-operation-server/pkg/ws-proxy/schema"
)

var (
	sceneClient pb.SceneServiceClient
	trackClient pb.TrackServiceClient
	debugClient pb.DebugServiceClient
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
	debugClient = pb.NewDebugServiceClient(conn)
	return nil
}
