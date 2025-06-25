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

func AddInternalRouterForTopicSolutionController(
	r *gin.RouterGroup,
	s app.TopicSolutionAppService,
) {
	ctl := TopicSolutionController{
		appService: s,
	}

	r.POST("/v1/hot-topic/:community/solution", ctl.Add)
}

type TopicSolutionController struct {
	appService app.TopicSolutionAppService
}

// @Summary      ToReview
// @Description  add topic solution
// @Tags         HotTopic
// @Param        community   path    string             true    "lowercase community name, like openubmc, cann"
// @Param        body        body    reqToAddSolution   true    "body"
// @Accept       json
// @Security     Internal
// @Success      201    {object}    commonctl.ResponseData{}
// @Router       /v1/hot-topic/{community}/solution [post]
func (ctl *TopicSolutionController) Add(ctx *gin.Context) {
	req := reqToAddSolution{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	if err := ctl.appService.Add(ctx.Param("community"), cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, nil)
	}
}
