/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"github.com/gin-gonic/gin"

	"github.com/opensourceways/hot-topic-website-backend/common/infrastructure/mongodb"
	"github.com/opensourceways/hot-topic-website-backend/config"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/app"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/controller"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/infrastructure/repositoryimpl"
)

func initHotTopicManagement(cfg *config.Config, services *allServices) error {
	htCfg := &cfg.HotTopicManagement
	hm := map[string]repositoryimpl.Dao{}
	nhm := map[string]repositoryimpl.Dao{}

	items := htCfg.Repo.CommunityCollections
	for i := range items {
		item := &items[i]
		hm[item.Community] = mongodb.DAO(item.Collections.HotTopic)
		nhm[item.Community] = mongodb.DAO(item.Collections.NotHotTopic)
	}

	services.hottopicmanagementApp = app.NewAppService(
		&htCfg.App,
		repositoryimpl.NewHotTopic(hm),
		repositoryimpl.NewNotHotTopic(nhm),
	)

	services.topicSolutionApp = app.NewTopicSolutionAppService(
		repositoryimpl.NewTopicSolution(mongodb.DAO(htCfg.Repo.TopicSolution)),
		repositoryimpl.NewHotTopic(hm),
	)

	return nil
}

func setInternalRouterForHotTopicManagement(rg *gin.RouterGroup, services *allServices) {
	controller.AddInternalRouterForHotTopicController(
		rg,
		services.hottopicmanagementApp,
	)
}

func setInternalRouterForTopicSolution(rg *gin.RouterGroup, services *allServices) {
	controller.AddInternalRouterForTopicSolutionController(
		rg,
		services.topicSolutionApp,
	)
}
