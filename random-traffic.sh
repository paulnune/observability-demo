#!/usr/bin/env bash

# Tempo máximo em segundos (5 min = 300s)
DURATION=$((5 * 60))
START=$(date +%s)

# Função para rodar um comando a cada segundo
run_loop() {
  local CMD=$1
  while true; do
    NOW=$(date +%s)
    ELAPSED=$((NOW - START))
    if [ $ELAPSED -ge $DURATION ]; then
      echo "[INFO] Finalizando loop do comando: $CMD"
      break
    fi
    echo "[INFO] Executando: $CMD"
    eval "$CMD" >/dev/null 2>&1
    sleep 1
  done
}

# Comandos
CMD1="curl -s -X POST http://localhost:8081/generate-log"
CMD2="curl -s -X POST http://localhost:8080/order"

# Executa em paralelo
run_loop "$CMD1" &
run_loop "$CMD2" &

# Aguarda todos terminarem
wait
echo "[INFO] Execução concluída após $DURATION segundos."
