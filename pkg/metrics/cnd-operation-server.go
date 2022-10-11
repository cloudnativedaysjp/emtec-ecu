package metrics

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/xerrors"
)

func RunCndOperationServer(addr string) error {
	m := echo.New()
	reg := prometheus.NewRegistry()
	RegisterCndOperationServer(reg)
	m.GET("/metrics",
		echo.WrapHandler(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))
	if err := m.Start(addr); err != nil {
		return xerrors.Errorf("message: %w", err)
	}
	return nil
}

func RegisterCndOperationServer(registry prometheus.Registerer) {
	const subsystem = "server"
	registerDreamkast(registry, subsystem)
}
