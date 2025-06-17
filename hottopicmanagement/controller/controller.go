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

func AddInternalRouterForHotTopicController(
	r *gin.RouterGroup,
	s app.AppService,
) {
	ctl := HotTopicController{
		appService: s,
	}

	r.POST("/v1/hot-topic/to-review", ctl.ToReview)
}

type HotTopicController struct {
	appService app.AppService
}

// @Summary      ToReview
// @Description  upload topics to review
// @Tags         HotTopic
// @Param        community   path    string        true    "lowercase community name, like openubmc, cann"
// @Param        body        body    reqToReview   true    "body"
// @Accept       json
// @Security     Internal
// @Success      201    {object}    commonctl.ResponseData{}
// @Router       /v1/hot-topic/{community}/to-review [post]
func (ctl *HotTopicController) ToReview(ctx *gin.Context) {
	req := reqToReview{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	if err := ctl.appService.ToReview(ctx.Param("community"), cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, nil)
	}
}
