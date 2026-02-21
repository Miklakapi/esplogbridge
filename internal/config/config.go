package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func Load(path string) (Config, error) {
	var cfg Config

	data, err := loadBytes(path)
	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	normalize(&cfg)
	applyDefaults(&cfg)

	if err := validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func loadBytes(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func normalize(cfg *Config) {
	cfg.Listen = strings.TrimSpace(cfg.Listen)

	cfg.Loki.URL = strings.TrimSpace(cfg.Loki.URL)
	cfg.Loki.Labels.Job = strings.TrimSpace(cfg.Loki.Labels.Job)

	if cfg.Loki.Headers != nil {
		clean := make(map[string]string, len(cfg.Loki.Headers))
		for k, v := range cfg.Loki.Headers {
			k = strings.TrimSpace(k)
			v = strings.TrimSpace(v)
			if k == "" || v == "" {
				continue
			}
			clean[k] = v
		}
		cfg.Loki.Headers = clean
	}

	if cfg.Devices != nil {
		clean := make(map[string]string, len(cfg.Devices))
		for ip, id := range cfg.Devices {
			ip = strings.TrimSpace(ip)
			id = strings.TrimSpace(id)
			if ip == "" || id == "" {
				continue
			}
			clean[ip] = id
		}
		cfg.Devices = clean
	}
}

func applyDefaults(cfg *Config) {
	if cfg.Input.UDP.QueueSize < 1 {
		cfg.Input.UDP.QueueSize = 100
	}

	if cfg.Loki.Timeout <= 0 {
		cfg.Loki.Timeout = 2 * time.Second
	}

	if cfg.Loki.Labels.Job == "" {
		cfg.Loki.Labels.Job = "esphome"
	}

	if cfg.Loki.Batch.MaxItems < 1 {
		cfg.Loki.Batch.MaxItems = 50
	}
	if cfg.Loki.Batch.MaxWait <= 0 {
		cfg.Loki.Batch.MaxWait = 1 * time.Second
	}
}

func validate(cfg Config) error {
	if cfg.Listen == "" {
		return errors.New("config: listen is required (e.g. :5514)")
	}

	if cfg.Input.UDP.QueueSize < 1 {
		return fmt.Errorf("config: input.udp.queue_size must be >= 1 (got: %d)", cfg.Input.UDP.QueueSize)
	}

	if cfg.Loki.URL == "" {
		return errors.New("config: loki.url is required")
	}

	if cfg.Loki.Timeout <= 0 {
		return errors.New("config: loki.timeout must be > 0")
	}

	if cfg.Loki.Labels.Job == "" {
		return errors.New("config: loki.labels.job cannot be empty")
	}

	if cfg.Loki.Batch.MaxItems < 1 {
		return fmt.Errorf("config: loki.batch.max_items must be >= 1 (got: %d)", cfg.Loki.Batch.MaxItems)
	}

	if cfg.Loki.Batch.MaxWait <= 0 {
		return errors.New("config: loki.batch.max_wait must be > 0")
	}

	if len(cfg.Devices) == 0 {
		return errors.New("config: devices must contain at least one IP -> device_id mapping")
	}

	for ip, id := range cfg.Devices {
		if ip == "" || id == "" {
			return errors.New("config: devices contains empty ip or device_id")
		}
	}

	return nil
}
