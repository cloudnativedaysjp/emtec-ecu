package server

import (
	"context"
	"net"

	"github.com/go-logr/logr"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/cloudnativedaysjp/emtec-ecu/pkg/infra/obsws"
	"github.com/cloudnativedaysjp/emtec-ecu/pkg/infra/sharedmem"
	pb "github.com/cloudnativedaysjp/emtec-ecu/pkg/ws-proxy/schema"
)

const componentName = "ws-proxy"

type Config struct {
	Development bool
	Logger      logr.Logger
	ZapLogger   *zap.Logger
	BindAddr    string
	Obs         []ConfigObs
}

type ConfigObs struct {
	Host      string
	Password  string
	DkTrackId int32
}

func Run(ctx context.Context, conf Config) error {
	logger := conf.Logger.WithName(componentName)

	// TODO(#57): move to cmd/server/main.go
	obswsClientMap := make(map[int32]obsws.Client)
	for _, obs := range conf.Obs {
		obswsClient, err := obsws.NewObsWebSocketClient(obs.Host, obs.Password)
		if err != nil {
			err := xerrors.Errorf("message: %w", err)
			logger.Error(err, "obsws.NewObsWebSocketClient() was failed")
			return err
		}
		obswsClientMap[obs.DkTrackId] = obswsClient
	}

	controller := &Controller{
		Logger:      logger,
		ObsWsMap:    obswsClientMap,
		MemWriter:   sharedmem.Writer{UseStorageForDisableAutomation: true},
		MemReader:   sharedmem.Reader{UseStorageForTrack: true, UseStorageForDisableAutomation: true},
		MemDebugger: sharedmem.Debugger{},
	}

	// Initialize gRPC server
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(conf.ZapLogger.Named(componentName)),
		),
	)
	pb.RegisterSceneServiceServer(s, controller)
	pb.RegisterTrackServiceServer(s, controller)
	pb.RegisterDebugServiceServer(s, controller)
	if conf.Development {
		reflection.Register(s)
	}

	// Serve
	serverErrStream := make(chan error)
	{
		lis, err := net.Listen("tcp", conf.BindAddr)
		if err != nil {
			err := xerrors.Errorf("message: %w", err)
			logger.Error(err, "net.Listen() was failed")
			return err
		}
		go func() {
			if err := s.Serve(lis); err != nil {
				serverErrStream <- err
			}
		}()
	}
	select {
	case <-ctx.Done():
		logger.Info("context was done.")
		return nil
	case err := <-serverErrStream:
		logger.Error(err, "s.Serve() was failed")
		return err
	}
}
