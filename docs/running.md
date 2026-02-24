# Running esplogbridge

After installation and configuration, esplogbridge can be run in two modes:
manual (foreground) or as a systemd service.

## Foreground mode (manual)

Run the agent:

```bash
esplogbridge --config esplogbridge.yaml
```

Stop the bridge with `Ctrl+C`.

On shutdown, the process:

- stops receiving UDP packets,
- flushes any remaining batch to Loki (best-effort),
- exits with non-zero code on fatal errors.

## Production usage (systemd service)

Running esplogbridge as a systemd service is recommended for long-running setups.

An example unit file is provided as [`example.service`](../example.service).<br>
Adapt paths to your environment:

- binary path (e.g. /usr/local/bin/esplogbridge),

- configuration path (e.g. /etc/esplogbridge.yaml).

Typical installation:

```bash
sudo cp example.service /etc/systemd/system/esplogbridge.service
sudo systemctl daemon-reload
sudo systemctl enable --now esplogbridge
```

Check status and logs:

```bash
systemctl status esplogbridge
journalctl -u esplogbridge -f
```
