package repository

import "github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"

type TopicSolutions struct {
	Id             string
	RetryNum       int
	Community      string
	TopicSolutions []domain.TopicSolution
}

type RepoTopicSolution interface {
	Add(string, []domain.TopicSolution) error
	Save(*TopicSolutions) error
	FindOldest() (TopicSolutions, error)
	Remove(string) error
}
