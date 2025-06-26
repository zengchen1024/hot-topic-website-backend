package repositoryimpl

import (
	"time"

	commonRepo "github.com/opensourceways/hot-topic-website-backend/common/domain/repository"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
	"go.mongodb.org/mongo-driver/bson"
)

func NewTopicReport(dao map[string]Dao) *topicReport {
	return &topicReport{
		daoMap: dao,
	}
}

type topicReport struct {
	daoMap
}

func (impl *topicReport) Add(community string, v []*domain.HotTopic) error {
	do := totopicReportDO(v)
	doc, err := do.toDoc()
	if err != nil {
		return err
	}

	docFilter := bson.M{fieldYear: doc[fieldYear], fieldWeek: doc[fieldWeek]}

	dao, err := impl.dao(community)
	if err != nil {
		return err
	}

	_, err = dao.InsertDocIfNotExists(docFilter, doc)
	if err != nil && dao.IsDocExists(err) {
		err = commonRepo.NewErrorDuplicateCreating(err)
	}

	return err
}

// insert a hot topic
func (impl *topicReport) InsertTopic(community string, year int, week int, v *domain.HotTopic) error {
	docFilter := bson.M{fieldYear: year, fieldWeek: week}

	insert := bson.M{"$set": bson.M{fieldOrder: v.Order}}
	currentVersion := time.Now().Day()

	dao, err := impl.dao(community)
	if err != nil {
		return err
	}
	err = dao.PushArraySingleItem(docFilter, fieldTopN, insert, currentVersion)
	if err != nil {
		err = commonRepo.NewErrorConcurrentUpdating(err)
	}
	update := bson.M{"$inc": bson.M{fieldCnt: 1}} // count +1

	err = dao.UpdateDocsWithoutVersion(docFilter, update)
	if err != nil {
		err = commonRepo.NewErrorConcurrentUpdating(err)
	}

	return err
}

// update a hottopic order index
func (impl *topicReport) UpdateTopic(community string, year int, week int, v *domain.HotTopic) error {
	docFilter := bson.M{fieldYear: year, fieldWeek: week}
	filterOfArray := bson.M{fieldTopic: v.Id}

	update := bson.M{"$set": bson.M{fieldOrder: v.Order}}
	currentVersion := time.Now().Day()

	dao, err := impl.dao(community)
	if err != nil {
		return err
	}
	err = dao.UpdateArraySingleItem(docFilter, fieldTopN, filterOfArray, update, currentVersion)
	if err != nil {
		err = commonRepo.NewErrorConcurrentUpdating(err)
	}

	return err
}

// Get all hot topic report
func (impl *topicReport) GetTopicReport(community string, year int, week int) (report domain.TopicReport, err error) {
	docFilter := bson.M{fieldYear: year, fieldWeek: week}

	dao, err := impl.dao(community)
	if err != nil {
		return report, err
	}

	sort := bson.M{
		fieldOrder: 1,
	}
	err = dao.GetDoc(docFilter, bson.M{}, sort, &report)
	if err != nil {
		err = commonRepo.NewErrorResourceNotFound(err)
	}

	return report, err
}

// Get all hot topic report
func (impl *topicReport) GetCurrentReport(community string) (report domain.TopicReport, err error) {
	year, week := time.Now().ISOWeek()
	return impl.GetTopicReport(community, year, week)
}

// Get  hot topic in last week
func (impl *topicReport) GetLastWeekTopic(community string) (report domain.TopicReport, err error) {
	now := time.Now()
	sevenDayAgo := now.AddDate(0, 0, -7)
	year, week := sevenDayAgo.ISOWeek()
	return impl.GetTopicReport(community, year, week)
}
