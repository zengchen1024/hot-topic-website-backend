package repositoryimpl

import (
	"go.mongodb.org/mongo-driver/bson"

	commonRepo "github.com/opensourceways/hot-topic-website-backend/common/domain/repository"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

func NewHotTopic(dao dao) *hotTopic {
	return &hotTopic{
		dao: dao,
	}
}

type hotTopic struct {
	dao dao
}

func (impl *hotTopic) Add(v *domain.HotTopic) error {
	do := tohotTopicDO(v)
	doc, err := do.toDoc()
	if err != nil {
		return err
	}
	doc[fieldVersion] = 0
	doc[fieldClosedAt] = 0

	docFilter := bson.M{fieldTitle: v.Title}

	_, err = impl.dao.InsertDocIfNotExists(docFilter, doc)
	if err != nil && impl.dao.IsDocExists(err) {
		err = commonRepo.NewErrorDuplicateCreating(err)
	}

	return err
}

func (impl *hotTopic) FindOpenOnes() ([]domain.HotTopic, error) {
	filter := bson.M{
		fieldClosedAt: 0,
	}

	sort := bson.M{
		fieldOrder: 1,
	}

	var dos []hotTopicDO

	if err := impl.dao.GetDocs(filter, nil, sort, &dos); err != nil {
		return nil, err
	}

	v := make([]domain.HotTopic, len(dos))
	for i := range dos {
		v[i] = dos[i].toHotTopic()
	}

	return v, nil
}
