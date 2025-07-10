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

func AddInternalRouterForTopicReviewController(
	r *gin.RouterGroup,
	s app.AppService,
) {
	ctl := TopicReviewController{
		appService: s,
	}

	r.POST("/v1/hot-topic/:community/to-review", ctl.Create)
	r.GET("/v1/topic-review/:community/publish", ctl.GetToPublish)
	r.GET("/v1/topic-review/:community", ctl.Get)
	r.PUT("/v1/topic-review/:community", ctl.Update)
}

type TopicReviewController struct {
	appService app.AppService
}

// @Summary      Create
// @Description  upload topics to review
// @Tags         TopicReview
// @Param        community   path    string        true    "lowercase community name, like openubmc, cann"
// @Param        body        body    reqToReview   true    "body"
// @Accept       json
// @Security     Internal
// @Success      201    {object}    commonctl.ResponseData{}
// @Router       /v1/hot-topic/{community}/to-review [post]
func (ctl *TopicReviewController) Create(ctx *gin.Context) {
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

	if err := ctl.appService.NewReviews(ctx.Param("community"), cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, nil)
	}
}

// @Summary      GetToPublish
// @Description  get topics to publish
// @Tags         TopicReview
// @Param        community   path    string        true    "lowercase community name, like openubmc, cann"
// @Accept       json
// @Security     Internal
// @Success      200    {object}    app.TopicsToPublishDTO{}
// @Router       /v1/topic-review/{community}/publish [get]
func (ctl *TopicReviewController) GetToPublish(ctx *gin.Context) {
	if v, err := ctl.appService.GetTopicsToPublish(ctx.Param("community")); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}

// @Summary      Get
// @Description  get topic review info
// @Tags         TopicReview
// @Param        community   path    string        true    "lowercase community name, like openubmc, cann"
// @Accept       json
// @Security     Internal
// @Success      200    {object}    app.TopicsToReviewDTO{}
// @Router       /v1/topic-review/{community} [get]
func (ctl *TopicReviewController) Get(ctx *gin.Context) {
	if v, err := ctl.appService.GetTopicsToReview(ctx.Param("community")); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}

// @Summary      Update
// @Description  update the selected topics
// @Tags         TopicReview
// @Param        community   path    string                true    "lowercase community name, like openubmc, cann"
// @Param        body        body    reqToUpdateSelected   true    "body"
// @Accept       json
// @Security     Internal
// @Success      202    {object}    commonctl.ResponseData{}
// @Router       /v1/topic-review/{community} [put]
func (ctl *TopicReviewController) Update(ctx *gin.Context) {
	req := reqToUpdateSelected{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	if err := req.Validate(); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	if err := ctl.appService.UpdateSelected(ctx.Param("community"), &req); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}
