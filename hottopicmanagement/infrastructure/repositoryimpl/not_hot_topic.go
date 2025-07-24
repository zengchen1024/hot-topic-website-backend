package repositoryimpl

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

func NewNotHotTopic(dao map[string]Dao) *notHotTopic {
	return &notHotTopic{
		daoMap: dao,
	}
}

type notHotTopic struct {
	daoMap
}

func (impl *notHotTopic) Save(community string, date int64, items []domain.NotHotTopic) error {
	dao, err := impl.dao(community)
	if err != nil {
		return err
	}

	if err := dao.DeleteDocs(bson.M{}); err != nil {
		return err
	}

	do := tonotHotTopicsDO(date, items)

	doc, err := do.toDoc()
	if err != nil {
		return err
	}

	_, err = dao.InsertDoc(doc)

	return err
}

func (impl *notHotTopic) FindAll(community string) ([]domain.NotHotTopic, error) {
	dao, err := impl.dao(community)
	if err != nil {
		return nil, err
	}

	var do notHotTopicsDO

	if err := dao.GetDoc(bson.M{}, nil, nil, &do); err != nil {
		return nil, err
	}

	return do.toNotHotTopics(), nil
}

func (impl *notHotTopic) FindCreatedAt(community string) (int64, error) {
	dao, err := impl.dao(community)
	if err != nil {
		return 0, err
	}

	var do notHotTopicsDO

	if err := dao.GetDoc(bson.M{}, bson.M{fieldCreatedAt: 1}, nil, &do); err != nil {
		if dao.IsDocNotExists(err) {
			return 0, nil
		}

		return 0, err
	}

	return do.CreatedAt, nil
}
