# prom-label-cleaner

> Utility to detect and prune high-cardinality labels from Prometheus metric exports

---

## Overview

`prom-label-cleaner` scans Prometheus `/metrics` endpoints or exported text files and identifies labels whose cardinality exceeds a configurable threshold. Offending labels can be stripped or replaced before the metrics are forwarded or stored.

---

## Installation

**Using `go install`:**

```bash
go install github.com/yourorg/prom-label-cleaner@latest
```

**Or build from source:**

```bash
git clone https://github.com/yourorg/prom-label-cleaner.git
cd prom-label-cleaner
go build -o prom-label-cleaner ./cmd/prom-label-cleaner
```

---

## Usage

```bash
# Scrape a local metrics endpoint and strip labels exceeding 100 unique values
prom-label-cleaner --source http://localhost:9090/metrics --threshold 100

# Process a static metrics file and write cleaned output
prom-label-cleaner --file ./metrics.txt --threshold 50 --output ./metrics_clean.txt

# Dry-run mode: report high-cardinality labels without modifying output
prom-label-cleaner --source http://localhost:9090/metrics --threshold 100 --dry-run
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--source` | — | URL of a Prometheus metrics endpoint |
| `--file` | — | Path to a local metrics text file |
| `--threshold` | `100` | Max unique label values before pruning |
| `--output` | stdout | File path for cleaned metrics output |
| `--dry-run` | `false` | Report issues without modifying output |

---

## License

This project is licensed under the [MIT License](LICENSE).