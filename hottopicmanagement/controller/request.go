package controller

import (
	"errors"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/app"
)

type reqToReview struct {
	Data []app.OptionalTopic `json:"data"`
}

func (req *reqToReview) toCmd() (app.CmdToUploadOptionalTopics, error) {
	if len(req.Data) == 0 {
		return nil, errors.New("no data")
	}

	return app.CmdToUploadOptionalTopics(req.Data), nil
}

type reqToAddSolution struct {
	Data []app.OptionalTopic `json:"data"`
}

func (req *reqToAddSolution) toCmd() (app.CmdToAddTopicSolution, error) {
	if len(req.Data) == 0 {
		return nil, errors.New("no data")
	}

	cmd := app.CmdToAddTopicSolution(req.Data)
	err := cmd.Validate()

	return cmd, err
}

type reqToUpdateSelected = app.CmdToUpdateSelected
