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

var timeLocation *time.Location

// LoadFromYaml reads a YAML file from the given path and unmarshals it into the provided interface.
func LoadFromYaml(path string, cfg interface{}) error {
	b, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, cfg)
}

func Now() time.Time {
	return time.Now().In(timeLocation)
}

func NowSec() int64 {
	return Now().Unix()
}

func Date() string {
	return Now().Format(layout)
}

func GetDate(t *time.Time) string {
	return t.Format(layout)
}

func DateToSecond(date string) int64 {
	t, err := time.Parse(layout, date)
	if err != nil {
		return 0
	}

	year, month, day := t.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, timeLocation).Unix()
}

func GetLastFriday() time.Time {
	t := Now()
	weekday := t.Weekday()

	var daysToSubtract int
	if weekday >= time.Friday {
		daysToSubtract = int(weekday - time.Friday)
	} else {
		daysToSubtract = int(weekday) + (7 - int(time.Friday))
	}

	t = t.AddDate(0, 0, -daysToSubtract)

	year, month, day := t.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func InitTimeZone() (err error) {
	timeLocation, err = time.LoadLocation("Asia/Shanghai")

	return
}
