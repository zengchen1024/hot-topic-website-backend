package repositoryimpl

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

func tonotNotHotTopicDO(v *domain.NotHotTopic) notNotHotTopicDO {
	return notNotHotTopicDO{
		Title:                 v.Title,
		Category:              v.Category,
		DiscussionSourceInfos: todiscussionSourceInfoDOs(v.DiscussionSources),
	}
}

func todiscussionSourceInfoDOs(items []domain.DiscussionSourceInfo) []discussionSourceInfoDO {
	r := make([]discussionSourceInfoDO, len(items))

	for i := range items {
		r[i] = todiscussionSourceInfoDO(&items[i])
	}

	return r
}

// notNotHotTopicDO
type notNotHotTopicDO struct {
	Title                 string                   `bson:"title"    json:"title"`
	Category              string                   `bson:"category" json:"category"`
	DiscussionSourceInfos []discussionSourceInfoDO `bson:"sources"  json:"sources"`
}

func (do *notNotHotTopicDO) toDoc() (bson.M, error) {
	return genDoc(do)
}

func (do *notNotHotTopicDO) toNotHotTopic() domain.NotHotTopic {
	return domain.NotHotTopic{
		Title:             do.Title,
		Category:          do.Category,
		DiscussionSources: do.toDiscussionSources(),
	}
}

func (do *notNotHotTopicDO) toDiscussionSources() []domain.DiscussionSourceInfo {
	r := make([]domain.DiscussionSourceInfo, len(do.DiscussionSourceInfos))

	for i := range do.DiscussionSourceInfos {
		r[i] = do.DiscussionSourceInfos[i].toDiscussionSourceInfo()
	}

	return r
}

// discussionSourceInfoDO
type discussionSourceInfoDO struct {
	Id     int    `bson:"id"     json:"id"`
	URL    string `bson:"url"    json:"url"`
	Title  string `bson:"title"  json:"title"`
	Closed bool   `bson:"closed" json:"closed"`
}

func (do *discussionSourceInfoDO) toDiscussionSourceInfo() domain.DiscussionSourceInfo {
	return domain.DiscussionSourceInfo{
		Id:     do.Id,
		URL:    do.URL,
		Title:  do.Title,
		Closed: do.Closed,
	}
}

func todiscussionSourceInfoDO(v *domain.DiscussionSourceInfo) discussionSourceInfoDO {
	return discussionSourceInfoDO{
		Id:     v.Id,
		URL:    v.URL,
		Title:  v.Title,
		Closed: v.Closed,
	}
}
