package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"

	"example.com/hello/cache"
)

const ExpireDurationStr = "1m"

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

var (
	inFlightGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_requests",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	})

	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"code", "method"},
	)

	// duration is partitioned by the HTTP method and handler. It uses custom
	// buckets based on the expected request duration.
	duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"handler", "method"},
	)
)

// Server as http server
type Server struct {
	cache cache.Cache
}

func (s *Server) start() {
	log.Info("Starting server on 8081 ...")

	// Register all of the metrics in the standard registry.
	prometheus.MustRegister(inFlightGauge, counter, duration)
	// Instrument the handlers with all the metrics, injecting the "handler"
	// label by currying.
	handlerChain := promhttp.InstrumentHandlerInFlight(inFlightGauge,
		promhttp.InstrumentHandlerDuration(duration.MustCurryWith(prometheus.Labels{"handler": "x"}),
			promhttp.InstrumentHandlerCounter(counter, validationMiddleware(s.handler)),
		),
	)

	mux := http.NewServeMux()
	mux.Handle("/", handlerChain)
	mux.HandleFunc("/size", s.size)
	mux.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8081", mux)
}

func (s *Server) handler(w http.ResponseWriter, req *http.Request) {
	key := req.URL.Path
	switch req.Method {
	case "GET":
		v := s.cache.Get(key)
		fmt.Fprintf(w, v)
	case "POST":
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Cannot read request body: %v", err), http.StatusBadRequest)
			return
		}
		value := string(body)
		s.cache.Set(key, value)
		fmt.Fprintf(w, value)

	default:
		http.Error(w, "Only GET and POST methods are supported", http.StatusMethodNotAllowed)
	}
}

func (s *Server) size(w http.ResponseWriter, req *http.Request) {
	size := s.cache.Size()
	fmt.Fprintf(w, "%d", size)
}

// Validate request url length and body length.
func validationMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 33 {
			http.Error(w, "Key length must be less than 32 char", http.StatusBadRequest)
			return
		}
		if len(r.URL.Path) < 2 {
			http.Error(w, "Must provide a key in the path", http.StatusBadRequest)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 32)
		f.ServeHTTP(w, r)
	}
}

func main() {
	expireDuration, _ := time.ParseDuration(ExpireDurationStr)
	cache := cache.NewSyncmapCache(expireDuration)
	server := Server{cache}
	server.start()
}
