# esplogbridge

![license](https://img.shields.io/badge/license-MIT-blue)
![linux](https://img.shields.io/badge/os-Linux-green)
![language](https://img.shields.io/badge/language-Go_1.25.1-blue)
![version](https://img.shields.io/github/v/tag/Miklakapi/esplogbridge)
![status](https://img.shields.io/badge/status-production-green)

A tiny UDP → Loki bridge for ESPHome logs.

**esplogbridge receives ESPHome UDP logs and forwards them to Grafana Loki using a single explicit YAML configuration file.**

The project focuses on:

- **one job only**: ESP UDP logs → Loki,
- **explicit configuration** (no hidden flags or modes),
- **bounded memory usage** (fixed queue, drop-oldest),
- **simple batching** for Loki Push API.

## Documentation

- [Installation](docs/installation.md)
- [Running esplogbridge](docs/running.md)

## Overview

esplogbridge is a small bridge designed to run on a Linux host and ship logs from ESP devices to Grafana Loki.

Key characteristics:

- single explicit YAML configuration file,
- UDP input with a bounded in-memory queue,
- lightweight IP-based device allowlist (`devices` map),
- Loki Push API output with simple per-device streams,
- deterministic batching (`max_items` / `max_wait`).

## Example configuration

A full example configuration is provided in [`example.yaml`](./example.yaml). Minimal version:

```yaml
listen: ':5514'

input:
    udp:
        queue_size: 100

loki:
    url: 'http://127.0.0.1:3100/loki/api/v1/push'
    timeout: 2s
    labels:
        job: 'esphome'
    batch:
        max_items: 50
        max_wait: 1s

devices:
    '192.168.1.55': 'esp-kitchen'
    '192.168.1.56': 'esp-bedroom'
```

## Technologies

Project is created with:

- Go 1.25.1
- systemd (optional)

## Status

The project's development has been completed.
