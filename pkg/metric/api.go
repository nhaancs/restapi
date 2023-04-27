package metric

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	asOnce  sync.Once
	apiServ *server
)

type server struct {
	reqTotal *prometheus.CounterVec
}

func Server() *server {
	asOnce.Do(func() {
		apiServ = &server{
			reqTotal: NewCounterVec(prometheus.CounterOpts{
				Name: "api_server_requests_total",
				Help: "The total number of handled API server requests.",
			}, []string{"method", "status", "code"}),
		}
		prometheus.MustRegister(apiServ)
	})

	return apiServ
}

func (s *server) Describe(desc chan<- *prometheus.Desc) {
	s.reqTotal.Describe(desc)
}

func (s *server) Collect(c chan<- prometheus.Metric) {
	s.reqTotal.Collect(c)
}

func (s *server) Inc(method, status, code string) {
	s.reqTotal.WithLabelValues(method, status, code).Inc()
}
