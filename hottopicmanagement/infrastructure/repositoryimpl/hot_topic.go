package repositoryimpl

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"

	commonRepo "github.com/opensourceways/hot-topic-website-backend/common/domain/repository"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

func NewHotTopic(dao map[string]Dao) *hotTopic {
	return &hotTopic{
		daoMap: dao,
	}
}

type hotTopic struct {
	daoMap
}

func (impl *hotTopic) Add(community string, v *domain.HotTopic) error {
	do := tohotTopicDO(v)
	doc, err := do.toDoc()
	if err != nil {
		return err
	}
	doc[fieldVersion] = 0
	doc[fieldClosedAt] = 0

	docFilter := bson.M{fieldTitle: v.Title}

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

func (impl *hotTopic) UpdateTopic(community string, v *domain.HotTopic) error {
	do := tohotTopicDO(v)
	doc, err := do.toDoc()
	if err != nil {
		return err
	}

	docFilter := bson.M{fieldId: v.Id}

	dao, err := impl.dao(community)
	if err != nil {
		return err
	}

	version := time.Now().Day()
	err = dao.UpdateDoc(docFilter, doc, version)

	return err
}

func (impl *hotTopic) FindOpenOnes(community string) ([]domain.HotTopic, error) {
	dao, err := impl.dao(community)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		fieldClosedAt: 0,
	}

	sort := bson.M{
		fieldOrder: 1,
	}

	var dos []hotTopicDO

	if err := dao.GetDocs(filter, nil, sort, &dos); err != nil {
		return nil, err
	}

	v := make([]domain.HotTopic, len(dos))
	for i := range dos {
		v[i] = dos[i].toHotTopic()
	}

	return v, nil
}

func (impl *hotTopic) FindOneWithId(community string, topic_id string) (topic domain.HotTopic, err error) {
	dao, err := impl.dao(community)
	if err != nil {
		return topic, err
	}

	filter := bson.M{fieldId: topic_id}

	sort := bson.M{
		fieldOrder: 1,
	}

	if err := dao.GetDoc(filter, nil, sort, &topic); err != nil {
		return topic, err
	}

	return topic, nil
}
