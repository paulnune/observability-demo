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
    └── docker-compose.yml  # Orquestração local (compatível com podman-compose)
```

---

## Pré-requisitos

- **Docker & Docker Compose** **ou** **Podman** (+ `podman-docker` e **`podman-compose`**)  
- Conta no **Dynatrace** com permissão para ingestão de logs  
- **.env** com variáveis do Dynatrace (veja abaixo)

> **Compatibilidade**  
> Esta demo é compatível com **Docker + Docker Compose**. Nos testes, foi executada em **RHEL 10** utilizando **Podman 5.4.0** com o pacote **`podman-docker`** (emulando o CLI `docker`) e **`podman-compose`**.

---

## 🚀 Preparação rápida

### 1) Clonar o repositório
```bash
git clone https://github.com/paulnune/observability-demo
cd observability-demo/deploy
```

### 2) Preparar variáveis do Dynatrace
```bash
cp .env.example .env
```
Edite o arquivo `.env` e preencha:
```bash
DT_ENV_URL="https://<tenant>.live.dynatrace.com"
DT_LOG_TOKEN="<api-token-com-log-ingest>"
```

### 3) Subir os serviços (escolha UMA das opções)

**Opção A: Docker**
```bash
docker compose up -d --build
```

**Opção B: Podman (recomendado em RHEL 10)**
```bash
podman-compose up -d --build
```

### 4) Verificar
```bash
curl -X POST http://localhost:8081/generate-log
curl -X POST http://localhost:8080/order
```

### 5) Testar aplicações
- **Legacy App** → <http://localhost:8081>  
  - `POST /generate-log` → gera um log não estruturado  
  - `GET /metrics` → expõe métricas  

- **Modern App** → <http://localhost:8080/order>  
  - `POST /order` → cria pedido e log estruturado  
  - `GET /metrics` → expõe métricas

### 6) Conferir ingestão no Dynatrace
No **Dynatrace → Grail → Logs**, filtre por:
```
dataset:demo
```
Ou via DQL:
```sql
fetch logs
| filter dataset == "demo"
```

---

## Exemplos de uso via curl

**Legacy App (logs não estruturados)**
```bash
curl -X POST http://localhost:8081/generate-log
```
Exemplo de resposta:
```json
{
  "log": "[2025-08-24 14:40:26] WARNING - Order failed due to out of stock (Order ID: 9661)",
  "message": "Log generated"
}
```

**Modern App (logs estruturados + OTEL)**
```bash
curl -X POST http://localhost:8080/order
```
Exemplo de resposta:
```
Order 802166 processed: 71 BRL
```

---

## Encerrar / resetar

**Parar serviços**
```bash
docker compose down            # Docker
# ou
podman-compose down            # Podman
```

**Remover volumes (reset total da demo)**
```bash
docker compose down -v         # Docker
# ou
podman-compose down -v         # Podman
```

---

## Cenários de Observabilidade demonstrados

1. **Normalização de logs não estruturados**  
   - Ex.: `"Payment failed for order"` → `severity=ERROR`, `loglevel=ERROR`.  

2. **Correlação de logs estruturados e traces**  
   - Modern App inclui `trace_id`, `span_id` e `service.name`.  

3. **Enriquecimento de contexto**  
   - Adição de campos como `dt.source.entity`, `dataset`, `service.name`.  

---

## Próximos passos

- Explorar ingestão direta via OTEL → Dynatrace Logs.  
- Incluir **dashboards comparativos** entre outras soluções, além do Dynatrace. 

---

## Licença

Este projeto está licenciado sob a [MIT License](./LICENSE).  