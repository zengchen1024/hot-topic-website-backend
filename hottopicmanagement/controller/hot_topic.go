/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	commonctl "github.com/opensourceways/hot-topic-website-backend/common/controller"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/app"
)

func AddInternalRouterForHotTopicController(
	r *gin.RouterGroup,
	s app.TopicSolutionAppService,
	s1 app.AppService,

) {
	ctl := HotTopicController{
		appService:    s,
		reviewService: s1,
	}

	r.POST("/v1/hot-topic/:community/solution", ctl.Add)
	r.GET("/v1/hot-topic/:community", ctl.Get)
}

type HotTopicController struct {
	appService    app.TopicSolutionAppService
	reviewService app.AppService
}

// @Summary      Add
// @Description  add topic solution
// @Tags         HotTopic
// @Param        community   path    string             true    "lowercase community name, like openubmc, cann"
// @Param        body        body    reqToAddSolution   true    "body"
// @Accept       json
// @Security     Internal
// @Success      201    {object}    commonctl.ResponseData{}
// @Router       /v1/hot-topic/{community}/solution [post]
func (ctl *HotTopicController) Add(ctx *gin.Context) {
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

// @Summary      Get
// @Description  get hot topics
// @Tags         HotTopic
// @Param        community   path    string             true    "lowercase community name, like openubmc, cann"
// @Param        since       query   int64              true    "get topics since the time(seconds)"
// @Accept       json
// @Security     Internal
// @Success      200    {object}    app.HotTopicsDTO{}
// @Router       /v1/hot-topic/{community} [get]
func (ctl *HotTopicController) Get(ctx *gin.Context) {
	var since int64
	if v, ok := ctx.GetQuery("since"); ok {
		v1, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			commonctl.SendBadRequestParam(ctx, err)

			return
		}

		since = v1
	}

	if v, err := ctl.reviewService.GetHotTopics(ctx.Param("community"), since); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}
