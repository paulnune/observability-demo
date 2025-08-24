# Overview of Observability Demo Project

The Observability Demo project is designed to showcase various applications and agents that generate and process logs and metrics for observability purposes. This project includes examples of both legacy and modern applications, as well as configurations for different observability tools.

## Project Structure

The project is organized into several key directories:

- **apps/**: Contains different applications that generate logs.
  - **legacy-app/**: A Python application that produces unstructured logs.
  - **modern-app/**: A Go application that generates structured logs in JSON format.
  - **ecommerce-java/**: A Java application simulating business metrics.

- **agents/**: Includes configurations for observability agents.
  - **vector/**: Configuration for Vector.dev, a tool for building observability pipelines.
  - **fluentbit/**: Contains configurations for Fluent Bit.
  - **otel-collector/**: Configuration for the OpenTelemetry Collector.

- **pipelines/**: Contains documentation and configurations for various observability pipelines.
  - **datadog/**: Pipeline for Datadog, handling logs and metrics.
  - **dynatrace/**: Configurations for Dynatrace.
  - **splunk/**: Configurations for Splunk Enterprise and Splunk Observability Cloud.
  - **grafana/**: Configurations for Grafana's Loki and Promtail.
  - **elastic/**: Configurations for Filebeat.

- **deploy/**: Contains deployment configurations and scripts.
  - **docker-compose.yml**: Orchestrates the deployment of applications and agents.
  - **.env.example**: Example environment variables for deployment.
  - **Makefile**: Commands for managing the demo.

- **docs/**: Documentation for the project, including an overview, comparisons, and architecture diagrams.

- **scripts/**: Contains scripts for generating logs for testing.

## Purpose

The main goal of this project is to demonstrate how different applications and agents can be integrated to provide observability through logging and metrics collection. By using a combination of legacy and modern applications, as well as various observability tools, users can gain insights into the performance and behavior of their systems.

This project serves as a practical example for developers and DevOps engineers looking to implement observability in their own applications and infrastructure.