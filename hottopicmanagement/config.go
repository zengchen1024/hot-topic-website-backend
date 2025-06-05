/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package models provides configuration and initialization functionality for the application.
package hottopicmanagement

import (
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/app"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/infrastructure/repositoryimpl"
)

// Config is a struct that represents the overall configuration for the application.
type Config struct {
	App  app.Config            `json:"app"`
	Repo repositoryimpl.Config `json:"repo"`
}

// ConfigItems returns a slice of interface{} containing pointers to the configuration items in the Config struct.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.App,
		&cfg.Repo,
	}
}
