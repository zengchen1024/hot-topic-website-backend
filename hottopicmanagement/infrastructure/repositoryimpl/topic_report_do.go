package repositoryimpl

import (
	"time"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	fieldYear    = "year"
	fieldWeek    = "week"
	fieldCnt     = "cnt"
	fieldUpateAt = "update_at"
	fieldTopN    = "topn"
	fieldTopic   = "topic"
)

func totopicReportDO(v []*domain.HotTopic) TopicReportDO {
	year, week := time.Now().ISOWeek()
	cnt := len(v)
	var topN []domain.TopicTopN
	var topicN domain.TopicTopN
	for _, topic := range v {
		topicN.TopicId = topic.Id
		topicN.Idx = topic.Order
		topN = append(topN, topicN)
	}

	return TopicReportDO{
		Year:     year,
		Week:     week,
		Cnt:      cnt,
		TopN:     topN,
		UpdateAt: time.Now().Unix(),
	}
}

type TopicReportDO struct {
	Year     int                `json:"year" bson:"year"`
	Week     int                `json:"week" bson:"year"`
	Cnt      int                `json:"cnt"  bson:"cnt"`
	TopN     []domain.TopicTopN `json:"topn" bson:"topn"`
	UpdateAt int64              `json:"_" bson:"update_at"`
}

func (do *TopicReportDO) toDoc() (bson.M, error) {
	return genDoc(do)
}
