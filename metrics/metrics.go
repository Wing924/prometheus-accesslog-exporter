package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "accesslog"

var labels = []string{"context"}

var (
	ReqTotal          = newCounterVec("requests", "total", "")
	ReqMethodTotal    = newCounterVec("requests", "method_total", "", "method")
	RespStatusTotal   = newCounterVec("response", "status_total", "", "code")
	RespCodeTotal     = newCounterVec("response", "code_total", "", "code")
	RespTimeHistogram *prometheus.HistogramVec
)

func newCounterVec(subsystem, name, help string, labels ...string) *prometheus.CounterVec {
	c := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      name,
		Help:      help,
		Subsystem: subsystem,
	}, append([]string{"context"}, labels...))
	prometheus.MustRegister(c)
	return c
}
