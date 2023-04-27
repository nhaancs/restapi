package metric

import (
	"database/sql"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func RegisterDB(db *sql.DB, name string) {
	if err := prometheus.Register(collectors.NewDBStatsCollector(db, name)); err != nil {
		panic(err)
	}
}

var (
	sOnce sync.Once
	s     *store
)

type store struct {
	reqTotal    *prometheus.CounterVec
	reqInFlight *prometheus.GaugeVec
	reqDur      *prometheus.HistogramVec
}

func Store() *store {
	sOnce.Do(func() {
		s = &store{
			reqTotal: NewCounterVec(prometheus.CounterOpts{
				Name: "storage_requests_total",
				Help: "The total number of handled storage requests.",
			}, []string{"name", "method"}),
			reqInFlight: NewGaugeVec(prometheus.GaugeOpts{
				Name: "storage_in_flight_requests",
				Help: "A gauge of requests currently being served.",
			}, []string{"name", "method"}),
			reqDur: NewHistogramVec(prometheus.HistogramOpts{
				Name: "storage_request_duration_seconds",
				Help: "A histogram of the request duration.",
			}, []string{"name", "method"}),
		}
		prometheus.MustRegister(s)
	})

	return s
}

func (s *store) Describe(desc chan<- *prometheus.Desc) {
	s.reqTotal.Describe(desc)
	s.reqInFlight.Describe(desc)
	s.reqDur.Describe(desc)
}

func (s *store) Collect(c chan<- prometheus.Metric) {
	s.reqTotal.Collect(c)
	s.reqInFlight.Collect(c)
	s.reqDur.Collect(c)
}

func (s *store) Start(name string, method string) func() {
	start := time.Now()
	s.reqInFlight.WithLabelValues(name, method).Inc()
	return func() {
		s.reqTotal.WithLabelValues(name, method).Inc()
		s.reqInFlight.WithLabelValues(name, method).Dec()
		s.reqDur.WithLabelValues(name, method).Observe(time.Since(start).Seconds())
	}
}
