# UptimeGo

Lightweight Go-based uptime monitoring CLI + daemon.  
Checks HTTP endpoints, exposes Prometheus metrics, and sends free alerts (Discord/Telegram).  
Simple, reliable, cloud-native observability tool.

## Features (MVP en cours)

- CLI intuitive (Cobra) : add/list/remove endpoints, start/stop daemon, status
- Periodic HTTP health checks (concurrent goroutines)
- Prometheus metrics endpoint (`/metrics`) : up/down gauge, latency histogram, counters
- Instant alerting via Discord webhook (ou Telegram/email en fallback)
- Configurable via YAML/TOML
- Graceful shutdown & retries with backoff
- Docker-ready

## Quick Start

### Prerequisites
- Go 1.22+
- Docker (optionnel pour démo Prometheus/Grafana)

### Installation (dev mode)

```bash
git clone https://github.com/tonusername/uptimego.git
cd uptimego
go mod tidy

Build & Run
Bash# Build the CLI
go build -o uptimego ./cmd/uptimego

# Create a basic config (see config.yaml.example)
cp config.yaml.example config.yaml

# Start monitoring
./uptimego start
Ou avec Docker (à venir) :
Bashdocker compose up
Configuration Example (config.yaml)
YAMLport: 8080
endpoints:
  - name: "My API Prod"
    url: "https://api.example.com/health"
    interval: 30s
    timeout: 5s
    down_threshold: 3
    alert:
      discord_webhook: "https://discord.com/api/webhooks/..."
Roadmap

 Project bootstrap & Cobra CLI
 Basic HTTP checker + concurrency
 Prometheus metrics instrumentation
 Alerting (Discord webhook)
 Docker + docker-compose with Prometheus/Grafana
 Tests & better README with screenshots

Why this project?
Built as a learning/portfolio project to demonstrate:

Go concurrency & error handling
Observability (Prometheus)
CLI tools in Go
Self-hosted monitoring basics

Perfect for junior backend/infra/cloud roles.
License
MIT License – see LICENSE

Made with ❤️ in Antananarivo, Madagascar
Questions or ideas? Open an issue!
