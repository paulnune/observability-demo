import random
from datetime import datetime
import os
from flask import Flask, jsonify

# Configurações
LOG_PATH = os.getenv("LOG_PATH", "/var/log/legacy-app/app.log")
os.makedirs(os.path.dirname(LOG_PATH), exist_ok=True)

# Contadores de métricas
metrics = {
    "orders_processed": 0,
    "payment_failures": 0,
    "out_of_stock_failures": 0,
    "total_errors": 0,
}

# Inicializa o Flask
app = Flask(__name__)

def generate_log():
    """Gera um log não estruturado com regras de negócio."""
    events = [
        {"type": "order_processed", "level": "INFO", "message": "Order processed successfully"},
        {"type": "payment_failure", "level": "ERROR", "message": "Payment failed for order"},
        {"type": "out_of_stock", "level": "WARNING", "message": "Order failed due to out of stock"},
    ]
    event = random.choice(events)
    now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    order_id = random.randint(1000, 9999)
    log_entry = f"[{now}] {event['level']} - {event['message']} (Order ID: {order_id})"

    # Atualiza métricas
    if event["type"] == "order_processed":
        metrics["orders_processed"] += 1
    elif event["type"] == "payment_failure":
        metrics["payment_failures"] += 1
        metrics["total_errors"] += 1
    elif event["type"] == "out_of_stock":
        metrics["out_of_stock_failures"] += 1
        metrics["total_errors"] += 1

    return log_entry

@app.route("/", methods=["GET"])
def home():
    """Endpoint raiz para informações básicas."""
    return jsonify({"message": "Legacy App is running", "endpoints": ["/generate-log", "/metrics"]})

@app.route("/generate-log", methods=["POST"])
def generate_log_endpoint():
    """Endpoint para gerar um log."""
    log_entry = generate_log()
    with open(LOG_PATH, "a") as log_file:
        log_file.write(log_entry + "\n")
    return jsonify({"message": "Log generated", "log": log_entry})

@app.route("/metrics", methods=["GET"])
def get_metrics():
    """Endpoint para expor métricas."""
    return jsonify(metrics)

if __name__ == "__main__":
    # Log de inicialização
    print(f"Legacy App is running. Logs will be written to {LOG_PATH}")
    print("Available endpoints: /generate-log (POST), /metrics (GET)")

    # Inicia o servidor Flask
    app.run(host="0.0.0.0", port=8081)