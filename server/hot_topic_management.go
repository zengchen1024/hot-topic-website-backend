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
	services.hottopicmanagementApp = app.NewAppService(
		&htCfg.App,
		repositoryimpl.NewHotTopic(mongodb.DAO(htCfg.Repo.Collections.HotTopic)),
		repositoryimpl.NewNotHotTopic(mongodb.DAO(htCfg.Repo.Collections.NotHotTopic)),
	)

	return nil
}

func setInternalRouterForHotTopicManagement(cfg *config.Config, rg *gin.RouterGroup, services *allServices) {
	controller.AddInternalRouterForHotTopicController(
		rg,
		services.hottopicmanagementApp,
	)
}
