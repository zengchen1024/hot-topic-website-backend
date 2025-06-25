/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"github.com/opensourceways/hot-topic-website-backend/config"
	hottopicmanagementapp "github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/app"
)

type allServices struct {
	topicSolutionApp      hottopicmanagementapp.TopicSolutionAppService
	hottopicmanagementApp hottopicmanagementapp.AppService
}

func initServices(cfg *config.Config) (services allServices, err error) {
	if err = initHotTopicManagement(cfg, &services); err != nil {
		return
	}

	return
}
