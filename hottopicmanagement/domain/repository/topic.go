package repository

import "github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"

type RepoHotTopic interface {
	Add(string, *domain.HotTopic) error
	FindOpenOnes(string) ([]domain.HotTopic, error)
	Find(string, string) (domain.HotTopic, error)
}

type RepoNotHotTopic interface {
	Add(string, *domain.NotHotTopic) error
	FindAll(string) ([]domain.NotHotTopic, error)
}
