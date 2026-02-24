package bridge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Miklakapi/esplogbridge/internal/config"
)

type Loki struct {
	cfg  config.Config
	http *http.Client
}

type lokiPush struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"`
}

func NewLoki(cfg config.Config) *Loki {
	return &Loki{
		cfg: cfg,
		http: &http.Client{
			Timeout: cfg.Loki.Timeout,
		},
	}
}

func (l *Loki) SendBatch(ctx context.Context, batch []Event) error {
	byDevice := make(map[string][]Event, 8)
	for _, ev := range batch {
		byDevice[ev.DeviceID] = append(byDevice[ev.DeviceID], ev)
	}

	deviceIDs := make([]string, 0, len(byDevice))
	for id := range byDevice {
		deviceIDs = append(deviceIDs, id)
	}

	streams := make([]lokiStream, 0, len(deviceIDs))

	for _, deviceID := range deviceIDs {
		events := byDevice[deviceID]

		values := make([][2]string, 0, len(events))
		for _, ev := range events {
			level := detectLevel(ev.RawLine)
			line := formatLine(ev.Timestamp, level, ev.RawLine)

			values = append(values, [2]string{
				fmt.Sprintf("%d", ev.Timestamp.UTC().UnixNano()),
				line,
			})
		}

		labels := map[string]string{
			"job":       l.cfg.Loki.Labels.Job,
			"device_id": deviceID,
		}

		streams = append(streams, lokiStream{
			Stream: labels,
			Values: values,
		})
	}

	body, err := json.Marshal(lokiPush{Streams: streams})
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, l.cfg.Loki.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range l.cfg.Loki.Headers {
		req.Header.Set(k, v)
	}

	resp, err := l.http.Do(req)
	if err != nil {
		return fmt.Errorf("post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("post: http %s", resp.Status)
	}

	return nil
}

func detectLevel(line string) string {
	if _, after, ok := strings.Cut(line, "level="); ok {
		rest := after
		if sp := strings.IndexByte(rest, ' '); sp >= 0 {
			rest = rest[:sp]
		}
		lvl := strings.ToLower(strings.TrimSpace(rest))
		switch lvl {
		case "debug", "info", "warn", "warning", "error":
			if lvl == "warning" {
				return "warn"
			}
			return lvl
		}
	}

	if len(line) >= 3 && line[0] == '[' && line[2] == ']' {
		switch line[1] {
		case 'D':
			return "debug"
		case 'I':
			return "info"
		case 'W':
			return "warn"
		case 'E':
			return "error"
		case 'C':
			return "info"
		}
	}

	return "info"
}

func formatLine(ts time.Time, level, msg string) string {
	return "ts=" + ts.UTC().Format(time.RFC3339Nano) +
		" level=" + level +
		" msg=" + escapeMsg(msg)
}

func escapeMsg(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return s
}
