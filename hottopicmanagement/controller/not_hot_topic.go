/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/opensourceways/hot-topic-website-backend/common/controller"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/app"
)

func AddInternalRouterForNotHotTopicController(
	r *gin.RouterGroup,
	s app.AppService,

) {
	ctl := NotHotTopicController{
		appService: s,
	}

	r.GET("/v1/not-hot-topic/:community", ctl.Get)
}

type NotHotTopicController struct {
	appService app.AppService
}

// @Summary      Get
// @Description  get worthless hot topics
// @Tags         NotHotTopic
// @Param        community   path    string    true    "lowercase community name, like openubmc, cann"
// @Accept       json
// @Security     Internal
// @Success      200    {object}    app.NotHotTopicsDTO{}
// @Router       /v1/not-hot-topic/{community} [get]
func (ctl *NotHotTopicController) Get(ctx *gin.Context) {
	if v, err := ctl.appService.GetWorthlessNotHotTopic(ctx.Param("community")); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}
