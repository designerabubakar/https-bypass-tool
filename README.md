# HTTPS Bypass Tool with TLS Fingerprint Spoofing

This Go-based tool is designed to perform HTTPS requests that bypass advanced anti-bot protections (e.g., Cloudflare, PerimeterX, Akamai) by spoofing TLS ClientHello fingerprints and using rotating SOCKS5 proxies.

## Features

- 🔒 Spoofed TLS fingerprints (Chrome, Firefox, Safari)
- 🧑‍💻 Browser-like HTTP headers
- 🕵️ Proxy rotation (SOCKS5 with authentication)
- ⏱️ Random delay between requests to mimic human behavior
- 📄 Easy to configure with `targets.txt` and `proxies.txt`

## Requirements

- Go 1.18+
- Modules:
  - [`utls`](https://github.com/refraction-networking/utls)
  - `golang.org/x/net/proxy`

## Installation

```bash
git clone https://github.com/your-username/https-bypass-tool.git
cd https-bypass-tool
go mod tidy
go build -o bypass https_bypass_tool.go
```

## Usage

Prepare two text files:

- `targets.txt`: List of URLs (one per line)
- `proxies.txt`: SOCKS5 proxies in the format `user:pass@host:port` (one per line)

Then run the tool:

```bash
./bypass
```

## Configuration

Default parameters (editable in `https_bypass_tool.go`):

- `numRequests` — Total number of requests (default: 50)
- `delayMin`, `delayMax` — Randomized delay between requests (in ms)
- `tlsProfiles` — Supported TLS fingerprints (`chrome`, `firefox`, `safari`)

## Legal Disclaimer

Usage against any website without explicit permission may violate their Terms of Service or applicable laws.

---
