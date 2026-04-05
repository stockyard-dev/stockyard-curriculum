# Stockyard Curriculum

**Self-hosted lesson planning and curriculum tracking**

Part of the [Stockyard](https://stockyard.dev) family of self-hosted tools.

## Quick Start

```bash
curl -fsSL https://stockyard.dev/tools/curriculum/install.sh | sh
```

Or with Docker:

```bash
docker run -p 9813:9813 -v curriculum_data:/data ghcr.io/stockyard-dev/stockyard-curriculum
```

Open `http://localhost:9813` in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9813` | HTTP port |
| `DATA_DIR` | `./curriculum-data` | SQLite database directory |
| `STOCKYARD_LICENSE_KEY` | *(empty)* | License key for unlimited use |

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Limits | 5 records | Unlimited |
| Price | Free | Included in bundle or $29.99/mo individual |

Get a license at [stockyard.dev](https://stockyard.dev).

## License

Apache 2.0
