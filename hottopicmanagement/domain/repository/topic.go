package repository

import "github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"

type RepoHotTopic interface {
	Add(string, *domain.HotTopic) error
	Find(string, string) (domain.HotTopic, error)
	Save(string, *domain.HotTopic) error
	FindAll(string) ([]domain.HotTopic, error)
	FindOpenOnes(string) ([]domain.HotTopic, error)
}

type RepoNotHotTopic interface {
	Save(community string, items []domain.NotHotTopic) error
	FindAll(string) ([]domain.NotHotTopic, error)
}

type RepoTopicsToReview interface {
	Add(string, *domain.TopicsToReview) error
	Find(string) (domain.TopicsToReview, error)
	SaveSelected(string, *domain.TopicsToReview) error
	FindSelected(string) (domain.TopicsToReview, error)
}
