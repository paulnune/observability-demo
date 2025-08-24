package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var (
	totalOrdersProcessed    int64
	totalErrors             int64
	totalProcessingTime     int64
	totalPaymentFailures    int64
	totalOutOfStockFailures int64

	serviceName string
)

// initTracer configura OTEL para enviar *traces* via OTLP (Collector).
func initTracer(ctx context.Context) (func(context.Context) error, error) {
	// Ex.: "http://otel-collector:4318"  (padrão)
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "http://otel-collector:4318"
	}

	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar exportador OTLP: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}

func handleOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tr := otel.Tracer(serviceName)
	_, span := tr.Start(ctx, "OrderProcessing")
	defer span.End()

	orderID := rand.Intn(1_000_000)
	amount := rand.Intn(100) + 1
	delay := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(delay)

	// Captura o contexto do span para correlação log↔trace
	sc := span.SpanContext()
	traceID := sc.TraceID().String()
	spanID := sc.SpanID().String()

	// Falha: out of stock (10%)
	if rand.Float32() < 0.10 {
		atomic.AddInt64(&totalErrors, 1)
		atomic.AddInt64(&totalOutOfStockFailures, 1)
		span.SetAttributes(attribute.String("order.status", "out_of_stock"))
		log.Printf("[WARN] Order %d failed: Out of stock trace_id=%s span_id=%s service.name=%s",
			orderID, traceID, spanID, serviceName)
		http.Error(w, "Out of stock", http.StatusConflict)
		return
	}

	// Falha: payment declined (15%)
	if rand.Float32() < 0.15 {
		atomic.AddInt64(&totalErrors, 1)
		atomic.AddInt64(&totalPaymentFailures, 1)
		span.SetAttributes(attribute.String("order.status", "payment_failed"))
		log.Printf("[ERROR] Order %d failed: Payment declined trace_id=%s span_id=%s service.name=%s",
			orderID, traceID, spanID, serviceName)
		http.Error(w, "Payment declined", http.StatusPaymentRequired)
		return
	}

	// Sucesso
	atomic.AddInt64(&totalOrdersProcessed, 1)
	atomic.AddInt64(&totalProcessingTime, int64(delay))

	span.SetAttributes(
		attribute.Int("order.id", orderID),
		attribute.Int("order.amount", amount),
		attribute.String("order.currency", "BRL"),
		attribute.String("order.status", "success"),
	)

	log.Printf("[INFO] Processed order %d: %d BRL trace_id=%s span_id=%s service.name=%s",
		orderID, amount, traceID, spanID, serviceName)
	fmt.Fprintf(w, "Order %d processed: %d BRL\n", orderID, amount)
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	totalOrders := atomic.LoadInt64(&totalOrdersProcessed)
	totalErr := atomic.LoadInt64(&totalErrors)
	totalTime := atomic.LoadInt64(&totalProcessingTime)
	paymentFailures := atomic.LoadInt64(&totalPaymentFailures)
	outOfStockFailures := atomic.LoadInt64(&totalOutOfStockFailures)

	avgProcessingTime := float64(0)
	if totalOrders > 0 {
		avgProcessingTime = float64(totalTime) / float64(totalOrders)
	}

	metrics := fmt.Sprintf(`# HELP orders_processed_total Total number of orders processed
# TYPE orders_processed_total counter
orders_processed_total %d

# HELP orders_failed_total Total number of failed orders
# TYPE orders_failed_total counter
orders_failed_total %d

# HELP orders_payment_failures_total Total number of payment failures
# TYPE orders_payment_failures_total counter
orders_payment_failures_total %d

# HELP orders_out_of_stock_failures_total Total number of out-of-stock failures
# TYPE orders_out_of_stock_failures_total counter
orders_out_of_stock_failures_total %d

# HELP orders_avg_processing_time Average processing time for orders
# TYPE orders_avg_processing_time gauge
orders_avg_processing_time %.2f
`, totalOrders, totalErr, paymentFailures, outOfStockFailures, avgProcessingTime)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(metrics))
}

func generatePeriodicLogs() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		totalOrders := atomic.LoadInt64(&totalOrdersProcessed)
		totalErr := atomic.LoadInt64(&totalErrors)
		paymentFailures := atomic.LoadInt64(&totalPaymentFailures)
		outOfStockFailures := atomic.LoadInt64(&totalOutOfStockFailures)
		// Periodic log informativo com service.name para facilitar filtros
		log.Printf("[INFO] Periodic log: %d orders processed, %d errors (Payment Failures: %d, Out of Stock: %d) service.name=%s",
			totalOrders, totalErr, paymentFailures, outOfStockFailures, serviceName)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Config por ENV
	serviceName = os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "modern-app"
	}
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = "/var/log/modern-app/app.log"
	}

	// Diretório + arquivo de log
	_ = os.MkdirAll("/var/log/modern-app", 0o755)
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("Erro ao abrir arquivo de log: %v", err)
	}
	defer f.Close()

	// Log em stdout + arquivo (Fluent Bit taila o arquivo)
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	// Mantém o timestamp padrão do logger do Go (compatível com o parser regex sugerido)

	// OTEL (traces)
	ctx := context.Background()
	shutdown, err := initTracer(ctx)
	if err != nil {
		log.Fatalf("Erro ao inicializar OTEL: %v", err)
	}
	defer func() {
		_ = shutdown(context.Background())
	}()

	// Handlers HTTP
	http.HandleFunc("/order", handleOrder)
	http.HandleFunc("/metrics", handleMetrics)

	// Log inicial
	log.Printf("[INFO] %s starting on :8080", serviceName)

	// Logs periódicos
	go generatePeriodicLogs()

	// Servidor
	log.Fatal(http.ListenAndServe(":8080", nil))
}
