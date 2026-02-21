package config

import "time"

type Config struct {
	Listen  string            `yaml:"listen"`
	Input   InputConfig       `yaml:"input"`
	Loki    LokiConfig        `yaml:"loki"`
	Devices map[string]string `yaml:"devices"`
}

type InputConfig struct {
	UDP UDPConfig `yaml:"udp"`
}

type UDPConfig struct {
	QueueSize int `yaml:"queue_size"`
}

type LokiConfig struct {
	URL     string            `yaml:"url"`
	Timeout time.Duration     `yaml:"timeout"`
	Labels  LokiLabels        `yaml:"labels"`
	Batch   BatchConfig       `yaml:"batch"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

type LokiLabels struct {
	Job string `yaml:"job"`
}

type BatchConfig struct {
	MaxItems int           `yaml:"max_items"`
	MaxWait  time.Duration `yaml:"max_wait"`
}
