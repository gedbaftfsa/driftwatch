# driftwatch

Detects config drift between deployed services and their declared infrastructure definitions.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git && cd driftwatch && go build ./...
```

---

## Usage

Point `driftwatch` at your infrastructure definition and a running environment to detect drift:

```bash
driftwatch scan --config ./infra/services.yaml --env production
```

Example output:

```
[DRIFT] service: api-gateway
  expected replicas: 3
  actual replicas:   2

[DRIFT] service: cache
  expected image: redis:7.0
  actual image:   redis:6.2

[OK] service: auth-service
```

### Common Flags

| Flag | Description |
|------|-------------|
| `--config` | Path to infrastructure definition file |
| `--env` | Target environment to scan |
| `--output` | Output format: `text`, `json`, `yaml` (default: `text`) |
| `--fail-on-drift` | Exit with code 1 if drift is detected (useful in CI) |

---

## Configuration

`driftwatch` supports YAML-based service definitions. See [`examples/`](./examples) for sample configs.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](./LICENSE)