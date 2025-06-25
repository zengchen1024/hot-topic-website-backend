package repository

import "github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"

type TopicSolutions struct {
	Id             string
	TopicSolutions []domain.TopicSolution
}

type RepoTopicSolution interface {
	Add(string, []domain.TopicSolution) error
	FindOldest() (TopicSolutions, error)
	Remove(string) error
}
