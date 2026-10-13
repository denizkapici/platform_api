package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Gelen HTTP istek sayısını takip eden Prometheus sayacı
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Toplam HTTP istek sayısı",
		},
		[]string{"path", "status"},
	)

	// İstek sürelerini ölçen histogram
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP istek süresi (saniye)",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// 1. Kubernetes Liveness Probe: Pod canlı mı?
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// 2. Kubernetes Readiness Probe: Pod trafik almaya hazır mı?
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		// İleride buraya DB bağlantı kontrolü (db.Ping) eklenecek
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	})

	// 3. Prometheus Metrik Ucu
	mux.Handle("/metrics", promhttp.Handler())

	// 4. Örnek API Ucu
	mux.HandleFunc("/api/v1/info", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			duration := time.Since(start).Seconds()
			httpRequestDuration.WithLabelValues("/api/v1/info").Observe(duration)
		}()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		httpRequestsTotal.WithLabelValues("/api/v1/info", "200").Inc()

		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Platform API çalışıyor",
			"version": "v1.0.0",
		})
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Graceful Shutdown dinleyicisi
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Sunucu :%s portunda dinleniyor...", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Sunucu hatası: %v", err)
		}
	}()

	// K8s SIGTERM sinyali gönderene kadar bekle
	<-stopChan
	log.Println("Kapatma sinyali alındı, sunucu güvenli şekilde durduruluyor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Zorla kapatma hatası: %v", err)
	}

	log.Println("Sunucu temiz bir şekilde kapatıldı.")
}