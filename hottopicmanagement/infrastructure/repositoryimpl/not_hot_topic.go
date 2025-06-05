package repositoryimpl

import (
	"go.mongodb.org/mongo-driver/bson"

	commonRepo "github.com/opensourceways/hot-topic-website-backend/common/domain/repository"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

func NewNotHotTopic(dao dao) *notNotHotTopic {
	return &notNotHotTopic{
		dao: dao,
	}
}

type notNotHotTopic struct {
	dao dao
}

func (impl *notNotHotTopic) Add(v *domain.NotHotTopic) error {
	do := tonotNotHotTopicDO(v)
	doc, err := do.toDoc()
	if err != nil {
		return err
	}

	docFilter := bson.M{fieldTitle: v.Title}

	_, err = impl.dao.InsertDocIfNotExists(docFilter, doc)
	if err != nil && impl.dao.IsDocExists(err) {
		err = commonRepo.NewErrorDuplicateCreating(err)
	}

	return err
}

func (impl *notNotHotTopic) FindAll() ([]domain.NotHotTopic, error) {
	var dos []notNotHotTopicDO

	if err := impl.dao.GetDocs(nil, nil, nil, &dos); err != nil {
		return nil, err
	}

	v := make([]domain.NotHotTopic, len(dos))
	for i := range dos {
		v[i] = dos[i].toNotHotTopic()
	}

	return v, nil
}
