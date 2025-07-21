package watchhottopic

import (
	"time"
)

type Config struct {
	// unit second
	Interval int `json:"interval"`
}

func (cfg *Config) SetDefault() {
	if cfg.Interval <= 0 {
		cfg.Interval = 24
	}
}

func (cfg *Config) intervalDuration() time.Duration {
	return time.Hour * time.Duration(cfg.Interval)
}
