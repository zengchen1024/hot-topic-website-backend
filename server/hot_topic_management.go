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
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/watch"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/watchhottopic"
)

func initHotTopicManagement(cfg *config.Config, services *allServices) error {
	htCfg := &cfg.HotTopicManagement
	hm := map[string]repositoryimpl.Dao{}
	nhm := map[string]repositoryimpl.Dao{}

	items := htCfg.Repo.CommunityCollections
	communities := make([]string, len(items))
	for i := range items {
		item := &items[i]
		hm[item.Community] = mongodb.DAO(item.Collections.HotTopic)
		nhm[item.Community] = mongodb.DAO(item.Collections.NotHotTopic)

		communities[i] = item.Community
	}

	services.repoHotTopic = repositoryimpl.NewHotTopic(hm)
	services.repoTopicSolution = repositoryimpl.NewTopicSolution(mongodb.DAO(htCfg.Repo.TopicSolution))
	repoNotHotTopic := repositoryimpl.NewNotHotTopic(nhm)

	services.hottopicmanagementApp = app.NewAppService(
		&htCfg.App,
		services.repoHotTopic,
		repoNotHotTopic,
		repositoryimpl.NewTopicToReview(mongodb.DAO(htCfg.Repo.TopicReview)),
	)

	services.topicSolutionApp = app.NewTopicSolutionAppService(
		services.repoTopicSolution, services.repoHotTopic,
	)

	watch.Start(&cfg.HotTopicManagement.Watch, services.repoHotTopic, services.repoTopicSolution)

	watchhottopic.Start(
		&cfg.HotTopicManagement.Apply, services.hottopicmanagementApp,
		repoNotHotTopic, communities,
	)

	return nil
}

func exitHotTopicManagement() {
	watch.Stop()
	watchhottopic.Stop()
}

func setInternalRouterForTopicReview(rg *gin.RouterGroup, services *allServices) {
	controller.AddInternalRouterForTopicReviewController(
		rg,
		services.hottopicmanagementApp,
	)
}

func setInternalRouterForNotHotTopic(rg *gin.RouterGroup, services *allServices) {
	controller.AddInternalRouterForNotHotTopicController(
		rg,
		services.hottopicmanagementApp,
	)
}

func setInternalRouterForHotTopic(rg *gin.RouterGroup, services *allServices) {
	controller.AddInternalRouterForHotTopicController(
		rg,
		services.topicSolutionApp,
		services.hottopicmanagementApp,
	)
}
