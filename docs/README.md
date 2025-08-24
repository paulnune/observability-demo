# Observability Demo Project

Welcome to the Observability Demo project! This repository showcases various applications and agents designed to demonstrate observability practices using different logging and monitoring tools.

## Project Structure

The project is organized into several directories, each serving a specific purpose:

- **apps/**: Contains different applications that generate logs.
  - **legacy-app/**: A legacy Python application that produces unstructured logs.
  - **modern-app/**: A modern Go application that generates structured logs in JSON format.
  - **ecommerce-java/**: A Java application simulating business metrics.

- **agents/**: Includes configurations for various observability agents.
  - **vector/**: Configuration for Vector.dev, a tool for building observability pipelines.
  - **fluentbit/**: Configurations for Fluent Bit.
  - **otel-collector/**: Configuration for the OpenTelemetry Collector.

- **pipelines/**: Contains documentation and configurations for different observability pipelines.
  - **datadog/**: Pipeline for Datadog handling logs and metrics.
  - **dynatrace/**: Configurations for Dynatrace.
  - **splunk/**: Configurations for Splunk Enterprise and Splunk Observability Cloud.
  - **grafana/**: Configurations for Loki and Promtail.
  - **elastic/**: Configurations for Filebeat.

- **deploy/**: Contains deployment configurations and scripts.
  - **docker-compose.yml**: Orchestrates the deployment of applications and agents.
  - **.env.example**: Example environment variables for deployment.
  - **Makefile**: Commands for running the demo.

- **.github/**: Contains GitHub Actions workflows for deployment.

- **docs/**: Documentation related to the project.
  - **overview.md**: Overview and explanation of the project.
  - **battlecard.md**: Comparison of observability tools.
  - **arch-diagram.drawio**: Architecture diagram.
  - **asciinema.cast**: Recording of the demo.

- **scripts/**: Contains scripts for generating logs.

## Getting Started

To get started with the Observability Demo project, follow these steps:

1. Clone the repository:
   ```
   git clone <repository-url>
   cd observability-demo
   ```

2. Set up the environment:
   - Copy the `.env.example` to `.env` and modify it as needed.

3. Build and run the applications and agents using Docker:
   ```
   docker-compose up --build
   ```

4. Explore the documentation in the `docs/` directory for more information on each component.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any suggestions or improvements.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.