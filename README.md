#  UptimeGo

[![Go Version](https://img.shields.io/github/go-mod/go-version/tonusername/uptimego)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**UptimeGo** is a lightweight, high-performance uptime monitoring CLI and daemon built with Go. It monitors your HTTP endpoints, exposes real-time Prometheus metrics, and sends instant alerts to Discord or Telegram when things go south.

Built with a "cloud-native first" mindset, it's perfect for developers who want a simple, self-hosted observability tool without the bloat.

---

##  Features (MVP in Progress)

- ** Intuitive CLI**: Built with [Cobra](https://github.com/spf13/cobra) for seamless management (`start`, `stop`, `status`, `list`).
- ** Concurrent Monitoring**: Checks multiple endpoints simultaneously using Go's powerful goroutines.
- ** Observability**: Native Prometheus metrics endpoint (`/metrics`) including:
  - Up/Down status gauges.
  - Response latency histograms.
  - Failure counters.
- ** Instant Alerting**: Discord Webhook support (Telegram & Email fallback coming soon).
- ** Flexible Config**: Easily manageable via `YAML` or `TOML`.
- ** Resilience**: Graceful shutdowns and automatic retries with backoff strategies.
- ** Docker Ready**: Containerized for easy deployment alongside your stack.

---

##  Quick Start

### Prerequisites

- **Go** 1.22+
- **Docker** (Optional, for Prometheus/Grafana demo)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Arnel-rah/uptimego.git
   cd uptimego
