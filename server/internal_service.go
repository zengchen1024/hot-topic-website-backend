/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"github.com/gin-gonic/gin"

	"github.com/opensourceways/hot-topic-website-backend/config"
)

func setInternalRouter(prefix string, engine *gin.Engine, cfg *config.Config, services *allServices) {
	rg := engine.Group(prefix)

	// set routers
	setInternalRouterForTopicReview(rg, services)

	setInternalRouterForTopicSolution(rg, services)
}
