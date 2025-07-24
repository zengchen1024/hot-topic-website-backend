/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"github.com/opensourceways/hot-topic-website-backend/config"
	hottopicmanagementapp "github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/app"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain/repository"
)

type allServices struct {
	topicSolutionApp      hottopicmanagementapp.TopicSolutionAppService
	hottopicmanagementApp hottopicmanagementapp.AppService

	repoHotTopic      repository.RepoHotTopic
	repoTopicSolution repository.RepoTopicSolution
}

func initServices(cfg *config.Config) (services allServices, err error) {
	if err = initHotTopicManagement(cfg, &services); err != nil {
		return
	}

	return
}

func exitServices() {
	exitHotTopicManagement()
}
