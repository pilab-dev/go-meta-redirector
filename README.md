# 🚀 GoMeta Redirector

[![CI/CD](https://github.com/pilab-dev/go-meta-redirector/actions/workflows/docker.yml/badge.svg)](https://github.com/pilab-dev/go-meta-redirector/actions/workflows/docker.yml)
[![Docker Pulls](https://img.shields.io/docker/pulls/ghcr.io/pilab-dev/go-meta-redirector)](https://github.com/pilab-dev/go-meta-redirector/pkgs/container/go-meta-redirector)
[![Go Report Card](https://goreportcard.com/badge/go.pilab.hu)](https://goreportcard.com/report/go.pilab.hu)
[![License](https://img.shields.io/github/license/pilab-dev/go-meta-redirector)](LICENSE)

> **Tiny, multi-domain Go module redirector** for custom `go get` domains. Supports `go.pilab.hu`, `go.paalgyula.com`, `go.pira.hu` with per-domain YAML configuration.

---

## ✨ Features

- 🌐 **Multi-Domain Support** - Route by `Host` header to different domain configs
- 📦 **Go Module Redirects** - Standard `go-import` meta tags for `go get`
- 📚 **Pkg.go.dev Integration** - Automatic redirects for documentation
- 🔄 **Fallback Patterns** - Wildcard matching for undefined repos
- 🐳 **Docker Ready** - Multi-stage build, minimal Alpine image
- ⚙️ **YAML Configuration** - Simple, human-readable config

---

## 🚀 Quick Start

### Using Docker (Recommended)

```bash
# Pull the latest image
docker pull ghcr.io/pilab-dev/go-meta-redirector:latest

# Run with default config
docker run -d -p 8080:8080 --name go-meta-redirector \
  ghcr.io/pilab-dev/go-meta-redirector:latest

# Run with custom config
docker run -d -p 8080:8080 -v $(pwd)/repos.yaml:/etc/go-meta-redirector/repos.yaml \
  ghcr.io/pilab-dev/go-meta-redirector:latest
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/pilab-dev/go-meta-redirector.git
cd go-meta-redirector

# Build binary
go build -o go-meta-redirector .

# Run (default port :8080, specify as argument)
./go-meta-redirector :8080
```

---

## 📝 Configuration

### `repos.yaml` Format

```yaml
domains:
  go.example.com:                    # Your Go module domain
    fallback:                         # Optional: fallback rule for undefined repos
      pattern: "cloud/*"              # Pattern with * wildcard
      target: "https://github.com/org/*"  # Target URL with * replacement
    repos:                            # Explicit repository mappings
      - path: cloud/hades             # Module path (after domain)
        git_url: https://github.com/org/hades.git
        pkgsite_url: https://pkg.go.dev/go.example.com/cloud/hades  # Optional

  go.another-domain.com:
    fallback: ~                        # Disable fallback (404 for undefined)
    repos: []
```

### Example: Pre-configured Domains

This repo comes pre-configured with:

| Domain | Fallback Pattern | Target |
|--------|-----------------|--------|
| `go.pilab.hu` | `cloud/*` | `github.com/pilab-dev/*` |
| `go.paalgyula.com` | `tools/*` | `github.com/paalgyula/*` |
| `go.pira.hu` | _(none)_ | _(explicit only)_ |

---

## 🔧 How It Works

### `go get` Flow

```bash
$ GOPROXY=direct go get go.pilab.hu/cloud/log
# 1. Go fetches: https://go.pilab.hu/cloud/log?go-get=1
# 2. Server returns meta tag:
#    <meta name="go-import" content="go.pilab.hu/cloud/log git https://github.com/pilab-dev/log.git">
# 3. Go downloads from GitHub
```

### Request Routing

```
┌─────────────────────┐
│  Incoming Request   │
│  Host: go.pilab.hu  │
│  Path: /cloud/log   │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Lookup Host in     │
│  repos.yaml         │
└──────────┬──────────┘
           │
           ▼
    ┌──────┴──────┐
    │              │
    ▼              ▼
┌─────────┐  ┌──────────┐
│ Exact   │  │ Fallback │
│ Match   │  │ Pattern  │
└────┬────┘  └────┬─────┘
     │              │
     ▼              ▼
┌─────────────────────────┐
│  Return go-import meta  │
│  or 404 if not found   │
└─────────────────────────┘
```

---

## 📦 Docker Images

Available tags on `ghcr.io/pilab-dev/go-meta-redirector`:

| Tag | Description |
|-----|-------------|
| `latest` | Latest main branch build |
| `main` | Main branch build |
| `v1.0.0` | Semantic version tag |
| `sha-xxxxxxx` | Commit SHA build |

### Custom Docker Build

```bash
# Build image
docker build -t go-meta-redirector .

# Run with custom port
docker run -p 8080:8080 go-meta-redirector :9090
```

---

## 🛠 Development

### Prerequisites

- Go 1.21+
- Docker (optional)

### Local Testing

```bash
# Start server
go run . :9090

# Test go-get response
curl -H "Host: go.pilab.hu" "http://localhost:9090/cloud/log?go-get=1"

# Test fallback
curl -H "Host: go.pilab.hu" "http://localhost:9090/cloud/unknown?go-get=1"

# Test browser redirect
curl -I -H "Host: go.pilab.hu" "http://localhost:9090/cloud/log"
```

### Project Structure

```
go-meta-redirector/
├── main.go          # Single-file server implementation
├── repos.yaml       # Multi-domain configuration
├── Dockerfile       # Multi-stage Docker build
├── go.mod           # Go module definition
└── .github/
    └── workflows/
        └── docker.yml   # CI/CD pipeline
```

---

## 🌍 Production Deployment

### DNS Setup

Point your Go module domains to your server:

```
go.pilab.hu    A    <your-server-ip>
go.paalgyula.com    A    <your-server-ip>
```

### Systemd Service

```ini
# /etc/systemd/system/go-meta-redirector.service
[Unit]
Description=GoMeta Redirector
After=network.target

[Service]
Type=simple
User=www-data
ExecStart=/usr/local/bin/go-meta-redirector :80
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### HTTPS with Let's Encrypt

```bash
# Use nginx/caddy as reverse proxy with TLS
# Or use cert-manager if running in Kubernetes
```

---

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

---

## 🤝 Contributing

Pull requests welcome! Please open an issue first for major changes.

---

<p align="center">
  Built with ❤️ by <a href="https://pilab.hu">Progressive Innovation LAB</a>
</p>
