package apm

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Business metrics
	ReservationCreated = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "business_reservation_created_total",
		Help: "Total number of reservations created",
	}, []string{"package_id", "destination"})

	ReservationFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "business_reservation_failed_total",
		Help: "Total number of failed reservation attempts",
	}, []string{"reason"})

	ReservationTravelersCaptured = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "business_reservation_travelers_captured_total",
		Help: "Total number of travelers captured",
	}, []string{"count"})

	ReservationExpired = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "business_reservation_expired_total",
		Help: "Total number of expired reservations",
	}, []string{"reason"})

	// Technical metrics
	APIDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "reservation_api_duration_seconds",
		Help:    "API request duration",
		Buckets: prometheus.DefBuckets,
	}, []string{"endpoint", "method"})

	DBDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "reservation_db_duration_seconds",
		Help:    "Database operation duration",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation"})

	CacheHit = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "reservation_cache_hit_total",
		Help: "Total cache hits",
	}, []string{"cache_name"})

	CacheMiss = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "reservation_cache_miss_total",
		Help: "Total cache misses",
	}, []string{"cache_name"})
)

// RecordReservationCreated registra criação de reserva
func RecordReservationCreated(packageID, destination string) {
	ReservationCreated.WithLabelValues(packageID, destination).Inc()
}

// RecordReservationFailed registra falha na criação
func RecordReservationFailed(reason string) {
	ReservationFailed.WithLabelValues(reason).Inc()
}

// RecordReservationTravelersCaptured registra captura de viajantes
func RecordReservationTravelersCaptured(count string) {
	ReservationTravelersCaptured.WithLabelValues(count).Inc()
}

// RecordReservationExpired registra expiração
func RecordReservationExpired(reason string) {
	ReservationExpired.WithLabelValues(reason).Inc()
}

// RecordAPIDuration registra duração de API
func RecordAPIDuration(endpoint, method string, duration time.Duration) {
	APIDuration.WithLabelValues(endpoint, method).Observe(duration.Seconds())
}

// RecordDBDuration registra duração de operação DB
func RecordDBDuration(operation string, duration time.Duration) {
	DBDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordCacheHit registra cache hit
func RecordCacheHit(cacheName string) {
	CacheHit.WithLabelValues(cacheName).Inc()
}

// RecordCacheMiss registra cache miss
func RecordCacheMiss(cacheName string) {
	CacheMiss.WithLabelValues(cacheName).Inc()
}

// ContextKey tipo para chaves de contexto
type ContextKey string

const (
	// TraceIDKey chave para trace ID no contexto
	TraceIDKey ContextKey = "trace_id"
	// SpanIDKey chave para span ID no contexto
	SpanIDKey ContextKey = "span_id"
)

// WithTraceID adiciona trace ID ao contexto
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// GetTraceID obtém trace ID do contexto
func GetTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(TraceIDKey).(string); ok {
		return id
	}
	return ""
}
