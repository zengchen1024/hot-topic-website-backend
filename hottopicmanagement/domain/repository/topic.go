package repository

import "github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"

type RepoHotTopic interface {
	Add(*domain.HotTopic) error
	FindOpenOnes() ([]domain.HotTopic, error)
}

type RepoNotHotTopic interface {
	Add(*domain.NotHotTopic) error
	FindAll() ([]domain.NotHotTopic, error)
}
