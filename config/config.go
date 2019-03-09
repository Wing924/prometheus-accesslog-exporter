package config

import (
	"time"

	"github.com/prometheus/common/model"
)

type (
	Config struct {
		ReadInterval model.Duration          `yaml:"read_interval"`
		Targets      map[string]TargetConfig `yaml:"targets"`
	}
	TargetConfig struct {
		Filepath  string            `yaml:"file_path"`
		Labels    map[string]string `yaml:"labels"`
		LogFormat string            `yaml:"log_format"`
		TimeScale float64           `yaml:"time_scale"`
	}
)

var (
	DefaultConfig = Config{
		ReadInterval: model.Duration(100 * time.Millisecond),
		Targets:      map[string]TargetConfig{},
	}
)
