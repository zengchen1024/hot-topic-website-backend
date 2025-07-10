package repositoryimpl

import (
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

func (impl *hotTopic) Save(community string, v *domain.HotTopic) error {
	do := tohotTopicDO(v)
	doc, err := do.toDoc()
	if err != nil {
		return err
	}

	dao, err := impl.dao(community)
	if err != nil {
		return err
	}

	docFilter, err := dao.DocIdFilter(v.Id)
	if err != nil {
		return err
	}

	return dao.UpdateDoc(docFilter, doc, v.Version)
}

func (impl *hotTopic) FindAll(community string) ([]domain.HotTopic, error) {
	dao, err := impl.dao(community)
	if err != nil {
		return nil, err
	}

	var dos []hotTopicDO

	if err := dao.GetDocs(bson.M{}, nil, nil, &dos); err != nil {
		return nil, err
	}

	v := make([]domain.HotTopic, len(dos))
	for i := range dos {
		v[i] = dos[i].toHotTopic()
	}

	return v, nil
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

func (impl *hotTopic) Find(community, topicId string) (topic domain.HotTopic, err error) {
	dao, err := impl.dao(community)
	if err != nil {
		return
	}

	filter, err := dao.DocIdFilter(topicId)
	if err != nil {
		return
	}

	var do hotTopicDO

	if err = dao.GetDoc(filter, nil, nil, &do); err != nil {
		return
	}

	topic = do.toHotTopic()

	return
}
