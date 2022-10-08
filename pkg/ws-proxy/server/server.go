package server

import (
	"context"
	"net"

	"github.com/go-logr/zapr"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/obsws"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/sharedmem"
	pb "github.com/cloudnativedaysjp/cnd-operation-server/pkg/ws-proxy/schema"
)

const componentName = "ws-proxy"

type Config struct {
	Development bool
	Debug       bool
	BindAddr    string
	Obs         []ConfigObs
}

type ConfigObs struct {
	Host      string
	Password  string
	DkTrackId int32
}

func Run(ctx context.Context, conf Config) error {
	// setup logger
	zapConf := zap.NewProductionConfig()
	if conf.Development {
		zapConf = zap.NewDevelopmentConfig()
	}

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
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(zapLogger.Named(componentName)),
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
