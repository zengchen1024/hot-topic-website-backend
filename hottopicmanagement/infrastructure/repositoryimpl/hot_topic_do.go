package repositoryimpl

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

const (
	fieldOrder     = "order"
	fieldTitle     = "title"
	fieldVersion   = "version"
	fieldClosedAt  = "closed_at"
	fieldCreatedAt = "created_at"
)

func tohotTopicDO(v *domain.HotTopic) hotTopicDO {
	return hotTopicDO{
		Title:             v.Title,
		Order:             v.Order,
		DiscussionSources: todiscussionSourceDOs(v.DiscussionSources),
		StatusTransferLog: tostatusLogDOs(v.StatusTransferLog),
	}
}

func todiscussionSourceDOs(items []domain.DiscussionSource) []discussionSourceDO {
	r := make([]discussionSourceDO, len(items))

	for i := range items {
		r[i] = todiscussionSourceDO(&items[i])
	}

	return r
}

func tostatusLogDOs(items []domain.StatusLog) []statusLogDO {
	r := make([]statusLogDO, len(items))

	for i := range items {
		r[i] = tostatusLogDO(&items[i])
	}

	return r
}

// hotTopicDO
type hotTopicDO struct {
	Id                primitive.ObjectID   `bson:"_id"            json:"-"`
	Title             string               `bson:"title"          json:"title"`
	Order             int                  `bson:"order"          json:"order"`
	DiscussionSources []discussionSourceDO `bson:"sources"        json:"sources"`
	StatusTransferLog []statusLogDO        `bson:"logs"           json:"logs"`
	Version           int                  `bson:"version"        json:"-"`
	ClosedAt          int64                `bson:"closed_at"      json:"-"`
	CreatedAt         int64                `bson:"created_at"     json:"created_at"`
}

func (do *hotTopicDO) toDoc() (bson.M, error) {
	return genDoc(do)
}

func (do *hotTopicDO) index() string {
	return do.Id.Hex()
}

func (do *hotTopicDO) toHotTopic() domain.HotTopic {
	return domain.HotTopic{
		Id:                do.index(),
		Title:             do.Title,
		Order:             do.Order,
		DiscussionSources: do.toDiscussionSources(),
		StatusTransferLog: do.toStatusLogs(),
		Version:           do.Version,
	}
}

func (do *hotTopicDO) toDiscussionSources() []domain.DiscussionSource {
	r := make([]domain.DiscussionSource, len(do.DiscussionSources))

	for i := range do.DiscussionSources {
		r[i] = do.DiscussionSources[i].toDiscussionSource()
	}

	return r
}

func (do *hotTopicDO) toStatusLogs() []domain.StatusLog {
	r := make([]domain.StatusLog, len(do.StatusTransferLog))

	for i := range do.StatusTransferLog {
		r[i] = do.StatusTransferLog[i].toStatusLog()
	}

	return r
}

// discussionSourceDO
type discussionSourceDO struct {
	Id         int    `bson:"id"           json:"id"`
	URL        string `bson:"url"          json:"url"`
	Type       string `bson:"type"         json:"type"`
	SourceId   string `bson:"source_id"    json:"source_id"`
	CreatedAt  string `bson:"created_at"   json:"created_at"`
	ImportedAt string `bson:"imported_at"  json:"imported_at"`
}

func (do *discussionSourceDO) toDiscussionSource() domain.DiscussionSource {
	return domain.DiscussionSource{
		DiscussionSourceMeta: domain.DiscussionSourceMeta{
			Id:        do.Id,
			URL:       do.URL,
			Type:      do.Type,
			SourceId:  do.SourceId,
			CreatedAt: do.CreatedAt,
		},
		ImportedAt: do.ImportedAt,
	}
}

func todiscussionSourceDO(v *domain.DiscussionSource) discussionSourceDO {
	return discussionSourceDO{
		Id:         v.Id,
		URL:        v.URL,
		Type:       v.Type,
		SourceId:   v.SourceId,
		CreatedAt:  v.CreatedAt,
		ImportedAt: v.ImportedAt,
	}
}

// statusLogDO
type statusLogDO struct {
	Time   string `bson:"time"   json:"time"`
	Status string `bson:"status" json:"status"`
}

func (do *statusLogDO) toStatusLog() domain.StatusLog {
	return domain.StatusLog{
		Time:   do.Time,
		Status: do.Status,
	}
}

func tostatusLogDO(v *domain.StatusLog) statusLogDO {
	return statusLogDO{
		Time:   v.Time,
		Status: v.Status,
	}
}
