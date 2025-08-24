package main

import (
    "context"
    "fmt"
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
    totalOrdersProcessed    int64 // Total de pedidos processados com sucesso
    totalErrors             int64 // Total de pedidos que falharam
    totalProcessingTime     int64 // Tempo total de processamento de pedidos
    totalPaymentFailures    int64 // Total de falhas de pagamento
    totalOutOfStockFailures int64 // Total de falhas por falta de estoque
)

func initTracer() func(context.Context) error {
    ctx := context.Background()

    exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint("otel-collector:4318"), otlptracehttp.WithInsecure())
    if err != nil {
        log.Fatalf("Erro ao criar exportador: %v", err)
    }

    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exp),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceName("modern-app"),
        )),
    )

    otel.SetTracerProvider(tp)
    return tp.Shutdown
}

func handleOrder(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    tr := otel.Tracer("modern-app")
    _, span := tr.Start(ctx, "OrderProcessing")
    defer span.End()

    // Gera um ID e valor aleatório para o pedido
    orderID := rand.Intn(1000000)
    amount := rand.Intn(100) + 1
    delay := time.Duration(rand.Intn(1000)) * time.Millisecond
    time.Sleep(delay)

    // Simula erro de falta de estoque em 10% dos pedidos
    if rand.Float32() < 0.1 {
        atomic.AddInt64(&totalErrors, 1)
        atomic.AddInt64(&totalOutOfStockFailures, 1)
        log.Printf("[WARN] Order %d failed: Out of stock", orderID)
        span.SetAttributes(attribute.String("order.status", "out_of_stock"))
        http.Error(w, "Out of stock", http.StatusConflict)
        return
    }

    // Simula erro de pagamento em 15% dos pedidos
    if rand.Float32() < 0.15 {
        atomic.AddInt64(&totalErrors, 1)
        atomic.AddInt64(&totalPaymentFailures, 1)
        log.Printf("[ERROR] Order %d failed: Payment declined", orderID)
        span.SetAttributes(attribute.String("order.status", "payment_failed"))
        http.Error(w, "Payment declined", http.StatusPaymentRequired)
        return
    }

    // Atualiza métricas de sucesso
    atomic.AddInt64(&totalOrdersProcessed, 1)
    atomic.AddInt64(&totalProcessingTime, int64(delay))

    // Adiciona atributos ao span
    span.SetAttributes(
        attribute.Int("order.id", orderID),
        attribute.Int("order.amount", amount),
        attribute.String("order.currency", "BRL"),
        attribute.String("order.status", "success"),
    )

    // Gera log de sucesso
    log.Printf("[INFO] Processed order %d: %d BRL", orderID, amount)
    fmt.Fprintf(w, "Order %d processed: %d BRL\n", orderID, amount)
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
    totalOrders := atomic.LoadInt64(&totalOrdersProcessed)
    totalErr := atomic.LoadInt64(&totalErrors)
    totalTime := atomic.LoadInt64(&totalProcessingTime)
    paymentFailures := atomic.LoadInt64(&totalPaymentFailures)
    outOfStockFailures := atomic.LoadInt64(&totalOutOfStockFailures)

    // Calcula o tempo médio de processamento
    avgProcessingTime := float64(0)
    if totalOrders > 0 {
        avgProcessingTime = float64(totalTime) / float64(totalOrders)
    }

    // Formata as métricas no formato Prometheus
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
    w.Write([]byte(metrics))
}

func generatePeriodicLogs() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        totalOrders := atomic.LoadInt64(&totalOrdersProcessed)
        totalErr := atomic.LoadInt64(&totalErrors)
        paymentFailures := atomic.LoadInt64(&totalPaymentFailures)
        outOfStockFailures := atomic.LoadInt64(&totalOutOfStockFailures)
        log.Printf("[INFO] Periodic log: %d orders processed, %d errors (Payment Failures: %d, Out of Stock: %d)",
            totalOrders, totalErr, paymentFailures, outOfStockFailures)
    }
}

func main() {
    log.SetOutput(os.Stdout)

    shutdown := initTracer()
    defer shutdown(context.Background())

    go generatePeriodicLogs()

    http.HandleFunc("/order", handleOrder)
    http.HandleFunc("/metrics", handleMetrics)

    log.Println("Modern app running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}