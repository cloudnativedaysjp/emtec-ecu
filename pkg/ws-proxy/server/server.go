package server

import (
	"fmt"
	"net"

	"github.com/go-logr/zapr"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/obsws"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/sharedmem"
	pb "github.com/cloudnativedaysjp/cnd-operation-server/pkg/ws-proxy/schema"
)

const componentName = "ws-proxy"

type Config struct {
	Debug    bool
	Obs      []ConfigObs
	BindAddr string
}

type ConfigObs struct {
	Host      string
	Password  string
	DkTrackId int32
}

func Run(conf Config) error {
	// setup logger
	zapConf := zap.NewProductionConfig()
	zapConf.DisableStacktrace = true // due to output wrapped error in errorVerbose
	zapLogger, err := zapConf.Build()
	if err != nil {
		return err
	}
	logger := zapr.NewLogger(zapLogger).WithName(componentName)

	obswsClientMap := make(map[int32]obsws.ClientIface)
	for _, obs := range conf.Obs {
		obswsClient, err := obsws.NewObsWebSocketClient(obs.Host, obs.Password)
		if err != nil {
			return err
		}
		obswsClientMap[obs.DkTrackId] = obswsClient
	}

	controller := &Controller{
		Logger:    logger,
		ObsWs:     obswsClientMap,
		Sharedmem: sharedmem.Writer{UseStorageForDisableAutomation: true},
	}

	// Initialize gRPC server
	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(zapLogger.Named(componentName)),
		),
	)
	pb.RegisterSceneServiceServer(s, controller)
	pb.RegisterTrackServiceServer(s, controller)
	if conf.Debug {
		reflection.Register(s)
	}

	// Serve
	lis, err := net.Listen("tcp", conf.BindAddr)
	if err != nil {
		return err
	}
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}
