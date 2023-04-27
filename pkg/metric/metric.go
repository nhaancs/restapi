package metric

import "github.com/prometheus/client_golang/prometheus"

const Namespace = "restapi"

func NewCounterVec(opts prometheus.CounterOpts, labelNames []string) *prometheus.CounterVec {
	opts.Namespace = Namespace
	return prometheus.NewCounterVec(opts, labelNames)
}

func NewGaugeVec(opts prometheus.GaugeOpts, labelNames []string) *prometheus.GaugeVec {
	opts.Namespace = Namespace
	return prometheus.NewGaugeVec(opts, labelNames)
}

func NewHistogramVec(opts prometheus.HistogramOpts, labelNames []string) *prometheus.HistogramVec {
	opts.Namespace = Namespace
	return prometheus.NewHistogramVec(opts, labelNames)
}
