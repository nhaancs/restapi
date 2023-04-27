package metric

import (
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	hsOnce sync.Once
	hs     *httpServer
)

type httpServer struct {
	reqTotal    *prometheus.CounterVec
	reqInFlight *prometheus.GaugeVec
	reqDur      *prometheus.HistogramVec
}

func HTTPServer() *httpServer {
	hsOnce.Do(func() {
		hs = &httpServer{
			reqTotal: NewCounterVec(prometheus.CounterOpts{
				Name: "http_server_requests_total",
				Help: "The total number of handled HTTP server requests.",
			}, []string{"method", "route", "status"}),
			reqInFlight: NewGaugeVec(prometheus.GaugeOpts{
				Name: "http_server_in_flight_requests",
				Help: "A gauge of requests currently being served by the server.",
			}, []string{"method", "route"}),
			reqDur: NewHistogramVec(prometheus.HistogramOpts{
				Name: "http_server_request_duration_seconds",
				Help: "A histogram of the request duration.",
			}, []string{"method", "route"}),
		}

		prometheus.MustRegister(hs)
	})

	return hs
}

func (hs *httpServer) Describe(desc chan<- *prometheus.Desc) {
	hs.reqTotal.Describe(desc)
	hs.reqInFlight.Describe(desc)
	hs.reqDur.Describe(desc)
}

func (hs *httpServer) Collect(c chan<- prometheus.Metric) {
	hs.reqTotal.Collect(c)
	hs.reqInFlight.Collect(c)
	hs.reqDur.Collect(c)
}

func (hs *httpServer) Start(method string, route string) func(status int) {
	start := time.Now()
	hs.reqInFlight.WithLabelValues(method, route).Inc()
	return func(status int) {
		hs.reqInFlight.WithLabelValues(method, route).Dec()
		hs.reqTotal.WithLabelValues(method, route, strconv.Itoa(status)).Inc()
		hs.reqDur.WithLabelValues(method, route).Observe(time.Since(start).Seconds())
	}
}
