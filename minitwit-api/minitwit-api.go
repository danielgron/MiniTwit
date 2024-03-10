package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"minitwit-api/api"

	"github.com/gorilla/mux"

	"minitwit-api/db"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
)

var (
	responseCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "minitwit_http_responses_total",
			Help: "The count of HTTP responses sent.",
		},
		[]string{"status"},
	)
	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "minitwit_request_duration_milliseconds",
			Help: "Request duration distribution.",
		},
		[]string{"path"},
	)
	cpuGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "minitwit_cpu_load_percent",
			Help: "Current CPU load as a percentage",
		},
	)
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w}

		cpuLoad, err := getCPULoad()
		if err != nil {
			log.Printf("Error getting CPU load: %v", err)
		} else {
			cpuGauge.Set(cpuLoad)
		}

		next.ServeHTTP(rw, r)

		responseCounter.WithLabelValues(strconv.Itoa(rw.status)).Inc()

		path, _ := mux.CurrentRoute(r).GetPathTemplate()
		timer := prometheus.NewTimer(requestDuration.WithLabelValues(path))
		defer timer.ObserveDuration()
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func getCPULoad() (float64, error) {
	percentages, err := cpu.Percent(0, false)
	if err != nil || len(percentages) == 0 {
		return 0, err
	}
	return percentages[0], nil
}

func main() {
	db.Connect_db()
	r := mux.NewRouter()

	r.Use(prometheusMiddleware)

	r.HandleFunc("/register", api.Register)
	r.HandleFunc("/msgs", api.Messages)
	r.HandleFunc("/msgs/{username}", api.Messages_per_user).Methods("GET", "POST")
	r.HandleFunc("/fllws/{username}", api.Follow)
	r.HandleFunc("/latest", api.Get_latest).Methods("GET")
	r.HandleFunc("/cleandb", api.Cleandb)
	r.HandleFunc("/delete", api.Delete)

	r.Handle("/metrics", promhttp.Handler())

	fmt.Println("Listening on port 15001...")
	err := http.ListenAndServe(":15001", r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
