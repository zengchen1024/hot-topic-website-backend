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
	closedAt := 0
	if v.IsResolved() {
		closedAt = 1
	}
	return hotTopicDO{
		Title:             v.Title,
		DiscussionSources: todiscussionSourceDOs(v.DiscussionSources),
		TransferLogs:      totransferLogDOs(v.TransferLogs),
		ClosedAt:          closedAt,
	}
}

func todiscussionSourceDOs(items []domain.DiscussionSource) []discussionSourceDO {
	r := make([]discussionSourceDO, len(items))

	for i := range items {
		r[i] = todiscussionSourceDO(&items[i])
	}

	return r
}

func totransferLogDOs(items []domain.TransferLog) []transferLogDO {
	r := make([]transferLogDO, len(items))

	for i := range items {
		r[i] = totransferLogDO(&items[i])
	}

	return r
}

// hotTopicDO
type hotTopicDO struct {
	Id                primitive.ObjectID   `bson:"_id"            json:"-"`
	Title             string               `bson:"title"          json:"title"`
	DiscussionSources []discussionSourceDO `bson:"sources"        json:"sources"`
	TransferLogs      []transferLogDO      `bson:"logs"           json:"logs"`
	ClosedAt          int                  `bson:"closed_at"      json:"closed_at"`
	Version           int                  `bson:"version"        json:"-"`
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
		DiscussionSources: do.toDiscussionSources(),
		TransferLogs:      do.toTransferLogs(),
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

func (do *hotTopicDO) toTransferLogs() []domain.TransferLog {
	r := make([]domain.TransferLog, len(do.TransferLogs))

	for i := range do.TransferLogs {
		r[i] = do.TransferLogs[i].toTransferLog()
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

// transferLogDO
type transferLogDO struct {
	Date   string `bson:"date"   json:"date"`
	Time   string `bson:"time"   json:"time"`
	Status string `bson:"status" json:"status"`
	Order  int    `bson:"order"  json:"order"`
}

func (do *transferLogDO) toTransferLog() domain.TransferLog {
	return domain.TransferLog{
		StatusLog: domain.StatusLog{
			Time:   do.Time,
			Status: do.Status,
		},
		Date:  do.Date,
		Order: do.Order,
	}
}

func totransferLogDO(v *domain.TransferLog) transferLogDO {
	return transferLogDO{
		Date:   v.Date,
		Time:   v.Time,
		Status: v.Status,
		Order:  v.Order,
	}
}
