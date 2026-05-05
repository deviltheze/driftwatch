# driftwatch

Lightweight daemon that detects config drift across remote servers via SSH.

---

## Installation

```bash
go install github.com/yourname/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/driftwatch.git && cd driftwatch && go build -o driftwatch .
```

---

## Usage

Define your servers and watched files in a `driftwatch.yaml` config file:

```yaml
interval: 60s
baseline: ./baseline
servers:
  - host: web-01.example.com
    user: deploy
    files:
      - /etc/nginx/nginx.conf
      - /etc/app/config.toml
  - host: web-02.example.com
    user: deploy
    files:
      - /etc/nginx/nginx.conf
```

Then run the daemon:

```bash
driftwatch --config driftwatch.yaml
```

driftwatch will SSH into each server on the defined interval, compare the target files against your local baseline snapshots, and report any differences:

```
[DRIFT] web-02.example.com:/etc/nginx/nginx.conf changed at 2024-11-03 14:22:01
[OK]    web-01.example.com:/etc/nginx/nginx.conf matches baseline
```

Use `--alert webhook` to post drift events to a Slack or custom HTTP endpoint.

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `driftwatch.yaml` | Path to config file |
| `--interval` | `60s` | Poll interval |
| `--alert` | `stdout` | Alert output (`stdout`, `webhook`) |
| `--webhook-url` | `""` | Webhook URL for alerts |

---

## License

MIT © yourname