# Observability Demo – Hybrid Logging & Business Events  

Este repositório contém uma Demo de Observabilidade. O objetivo é demonstrar como lidar com **logs estruturados e não estruturados**, enriquecendo e normalizando eventos antes de enviá-los para backends de observabilidade, com ênfase em **Dynatrace Grail**.  

A demo cobre **dois cenários de aplicações** (legado e moderno), agentes de coleta e normalização de logs (Fluent Bit + Lua), e pipeline de integração com OTEL Collector.  

Além disso, foram criados **Notebooks de Business Observability** no Dynatrace, disponíveis em dois formatos:  

- [Business Observability – Demo (DQL Input Visible)](https://szn23895.apps.dynatrace.com/ui/document/v0/#share=e89885fe-1849-4b22-878d-fa4d578d8aa7)  
  
- [Business Observability – Demo (DQL Input Hidden)](https://szn23895.apps.dynatrace.com/ui/document/v0/#share=54b0d423-4beb-4957-9715-376ca2c1cc1d)  

Também é possível visualizar as versões em PDF exportadas, localizadas no diretório [`files/`](./files):  

- [`BusinessObservability-Demo-Visible.pdf`](/files/Business%20Observability%20–%20Demo%20(DQL%20Input%20Visible).pdf)
  
- [`BusinessObservability-Demo-Hidden.pdf`](/files/Business%20Observability%20–%20Demo%20(DQL%20Input%20Hidden).pdf)

---

## Arquitetura

- **Legacy App (Python/Flask)**  
  - Gera **logs não estruturados** simulando regras de negócio (`processed order`, `payment declined`, `out of stock`).  
  - Exposição de métricas simples via `/metrics`.  

- **Modern App (Go + OpenTelemetry)**  
  - Gera **logs estruturados** já correlacionados com **traces OTLP**.  
  - Inclui atributos de negócio (ID do pedido, status, valor).  
  - Exposição de métricas Prometheus.  

- **Fluent Bit**  
  - Tail de arquivos de log das duas apps.  
  - Parsing com regex customizados (`parsers.conf`).  
  - Normalização de severidade (`normalize_severity.lua`).  
  - Classificação por mensagem para legado (`classify_by_message.lua`).  
  - Shape final compatível com Dynatrace (`to_dynatrace.lua`).  
  - Enriquecimento de campos: `service.name`, `dataset`, `dt.source.entity`.  
  - Envio via **HTTP → Dynatrace Logs v2 API**.  

- **OTEL Collector**  
  - Recebe logs via Fluent Forward.  
  - (Nesta demo: exportação debug, mas pronto para envio a SaaS como Dynatrace, Splunk, Grafana, Elastic, para efeito de comparação).  

---

## Estrutura

```
.
├── agents
│   ├── fluentbit        # Configuração Fluent Bit + Lua filters
│   └── otel-collector   # Configuração OTEL Collector
├── apps
│   ├── legacy-app       # Python Flask (logs não estruturados)
│   └── modern-app       # Go + OTEL (logs estruturados + traces)
└── deploy
    └── docker-compose.yml  # Orquestração local
```

---

## Como executar

Pré-requisitos:  
- Docker + Docker Compose  
- Variáveis de ambiente definidas no `.env` (Dynatrace tenant e token)  

### 1. Clonar repositório
```bash
git clone https://github.com/paulnune/observability-demo
cd observability-demo/deploy
```

### 2. Configurar secrets
Crie o arquivo `.env` com:
```bash
DT_ENV_URL="https://<tenant>.live.dynatrace.com"
DT_LOG_TOKEN="<api-token-com-log-ingest>"
```

### 3. Subir ambiente
```bash
docker compose up -d --build
```

### 4. Testar aplicações
- **Legacy App** → [http://localhost:8081](http://localhost:8081)  
  - `POST /generate-log` → gera um log não estruturado.  
  - `GET /metrics` → expõe métricas.  

- **Modern App** → [http://localhost:8080/order](http://localhost:8080/order)  
  - Gera pedido e log estruturado.  
  - Correlação log ↔ trace via OTEL.  
  - `GET /metrics` → expõe métricas.  

### 5. Verificar ingestão de logs
No **Dynatrace Grail → Logs**, filtre por `dataset:demo`.  

![alt text](/files/image.png)

---

## Exemplos de uso

Depois de subir o ambiente com `docker compose up -d --build`, é possível gerar logs diretamente via **curl**:

### Legacy App (logs não estruturados)
Gerar log manual:
```powershell
curl -X POST http://localhost:8081/generate-log
```
Exemplo de resposta:
```json
{
  "log": "[2025-08-24 14:40:26] WARNING - Order failed due to out of stock (Order ID: 9661)",
  "message": "Log generated"
}
```

Esse log será **parseado pelo Fluent Bit**, normalizado e classificado (`severity=WARN`, `loglevel=WARN`) antes de ser enviado ao Dynatrace.

---

### Modern App (logs estruturados + OTEL)
Criar um pedido:
```powershell
curl -X POST http://localhost:8080/order
```
Exemplo de resposta:
```
Order 802166 processed: 71 BRL
```

Esse evento gera um log estruturado já com `trace_id`, `span_id` e `service.name`, permitindo **correlação direta log ↔ trace**.

---

Para verificar no Dynatrace, filtre os logs em **Grail → Logs** usando:  
```sql
fetch logs
| filter dataset == "demo"
```

---

## Cenários de Observabilidade demonstrados

1. **Normalização de logs não estruturados**  
   - Exemplo: `"Payment failed for order"` → `severity=ERROR`, `loglevel=ERROR`.  

2. **Correlação de logs estruturados e traces**  
   - Modern App inclui `trace_id`, `span_id` e `service.name`.  

3. **Enriquecimento de contexto**  
   - Adição de campos como `dt.source.entity`, `dataset`, `service.name`.  

---

## 🔮 Próximos passos

- Expandir cenários para incluir **business observability** (KPIs de pedidos, falhas de pagamento etc. como logs/metrics).  
- Explorar ingestão direta via OTEL → Dynatrace Logs.  
- Incluir **dashboards comparativos** entre backends observáveis.

---

## 📜 Licença

Este projeto está licenciado sob a [MIT License](./LICENSE).  
