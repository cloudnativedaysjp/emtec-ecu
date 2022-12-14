package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//
// Metrics
//

var (
	dreamkastRequestSummaryVec prometheus.SummaryVec
)

func registerDreamkast(registry prometheus.Registerer, namespace string) {
	dreamkastRequestSummaryVec = *prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: namespace,
		Subsystem: "dreamkast",
		Name:      "request_second_total",
	}, []string{"endpointUrl", "kind"})
}

//
// Interface
//

type DreamkastMetricsIface interface {
	ListTracks(time.Duration)
	ListTalks(time.Duration)
	UpdateTalk(time.Duration)
}

//
// Data Access Object
//

type DreamkastMetricsDao struct {
	endpointUrl string
}

func NewDreamkastMetricsDao(endpointUrl string) *DreamkastMetricsDao {
	return &DreamkastMetricsDao{endpointUrl}
}

func (dao DreamkastMetricsDao) ListTracks(d time.Duration) {
	dreamkastRequestSummaryVec.
		WithLabelValues(dao.endpointUrl, "listTracks").Observe(float64(d / time.Second))
}

func (dao DreamkastMetricsDao) ListTalks(d time.Duration) {
	dreamkastRequestSummaryVec.
		WithLabelValues(dao.endpointUrl, "listTalks").Observe(float64(d / time.Second))
}

func (dao DreamkastMetricsDao) UpdateTalk(d time.Duration) {
	dreamkastRequestSummaryVec.
		WithLabelValues(dao.endpointUrl, "updateTalk").Observe(float64(d / time.Second))
}

//
// Fake Object
//

type DreamkastMetricsFake struct{}

func (DreamkastMetricsFake) ListTracks(time.Duration) {}
func (DreamkastMetricsFake) ListTalks(time.Duration)  {}
func (DreamkastMetricsFake) UpdateTalk(time.Duration) {}

//
// Utilities
//

type contextKeyDreamkastMetrics struct{}

var ctxKeyDreamkastMetrics = contextKeyDreamkastMetrics{}

func SetDreamkastMetricsToCtx(ctx context.Context, m DreamkastMetricsIface) context.Context {
	return context.WithValue(ctx, ctxKeyDreamkastMetrics, m)
}

func DreamkastMetricsFromCtx(ctx context.Context) DreamkastMetricsIface {
	dao, ok := ctx.Value(ctxKeyDreamkastMetrics).(DreamkastMetricsDao)
	if !ok {
		return &DreamkastMetricsFake{}
	}
	return dao
}
