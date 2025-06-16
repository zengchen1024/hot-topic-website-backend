/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package utils provides utility functions for various purposes.
package utils

import (
	"os"
	"time"

	"sigs.k8s.io/yaml"
)

const layout = "2006-01-02"

// LoadFromYaml reads a YAML file from the given path and unmarshals it into the provided interface.
func LoadFromYaml(path string, cfg interface{}) error {
	b, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, cfg)
}

func Now() int64 {
	return time.Now().Unix()
}

func Date() string {
	return time.Now().Format(layout)
}
