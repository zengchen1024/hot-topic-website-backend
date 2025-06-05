/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package config provides functionality for managing application configuration.
package config

import (
	"os"

	common "github.com/opensourceways/hot-topic-website-backend/common/config"
	"github.com/opensourceways/hot-topic-website-backend/common/infrastructure/mongodb"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement"
	"github.com/opensourceways/hot-topic-website-backend/utils"
)

// LoadConfig loads the configuration file from the specified path and deletes the file if needed
func LoadConfig(path string, cfg *Config, remove bool) error {
	if remove {
		defer os.Remove(path)
	}
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return err
	}

	common.SetDefault(cfg)

	return common.Validate(cfg)
}

type SwaggerInfo struct {
	Version string `json:"version" required:"true"`
	Title   string `json:"title" required:"true"`
	Desc    string `json:"desc" required:"true"`
}

// Config is a struct that represents the overall configuration for the application.
type Config struct {
	MongoDB            mongodb.Config            `json:"mongodb"`
	ReadHeaderTimeout  int                       `json:"read_header_timeout"`
	HotTopicManagement hottopicmanagement.Config `json:"hot_topic_management"`
}

// Init initializes the application using the configuration settings provided in the Config struct.
func (cfg *Config) Init() error {
	return nil
}

// ConfigItems returns a slice of interface{} containing pointers to the configuration items.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.HotTopicManagement,
	}
}

// SetDefault sets default values for the Config struct.
func (cfg *Config) SetDefault() {
	if cfg.ReadHeaderTimeout <= 0 {
		cfg.ReadHeaderTimeout = 10
	}
}

// Validate validates the configuration.
func (cfg *Config) Validate() error {
	return utils.CheckConfig(cfg, "")
}
