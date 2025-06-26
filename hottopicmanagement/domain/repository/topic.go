package repository

import "github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"

type RepoHotTopic interface {
	Add(string, *domain.HotTopic) error
	FindOpenOnes(string) ([]domain.HotTopic, error)
	FindOneWithId(string, string) (domain.HotTopic, error)
}

type RepoNotHotTopic interface {
	Add(string, *domain.NotHotTopic) error
	FindAll(string) ([]domain.NotHotTopic, error)
}

type RepoTopicReport interface {
	Add(string, []*domain.HotTopic) error
	InsertTopic(string, int, int, *domain.HotTopic) error
	UpdateTopic(string, int, int, *domain.HotTopic) error
	GetTopicReport(string, int, int) (domain.TopicReport, error)
	GetCurrentReport(string) (domain.TopicReport, error)
	GetLastWeekTopic(string) (domain.TopicReport, error)
}
