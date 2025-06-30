package repositoryimpl

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain/repository"
)

func NewTopicSolution(dao Dao) *topicSolution {
	return &topicSolution{
		dao: dao,
	}
}

type topicSolution struct {
	dao Dao
}

func (impl *topicSolution) Add(community string, v []domain.TopicSolution) error {
	do := totopicSolutionsDO(community, 0, v)
	doc, err := do.toDoc()
	if err != nil {
		return err
	}

	_, err = impl.dao.InsertDoc(doc)

	return err
}

func (impl *topicSolution) FindOldest() (repository.TopicSolutions, error) {
	sort := bson.M{
		fieldCreatedAt: 1,
	}

	var do topicSolutionsDO

	if err := impl.dao.GetDoc(bson.M{}, nil, sort, &do); err != nil {
		return repository.TopicSolutions{}, err
	}

	return do.toTopicSolutions(), nil
}

func (impl *topicSolution) Remove(tsId string) error {
	filter, err := impl.dao.DocIdFilter(tsId)
	if err != nil {
		return err
	}

	return impl.dao.DeleteDoc(filter)
}

func (impl *topicSolution) Save(solution *repository.TopicSolutions) error {
	do := totopicSolutionsDO(solution.Community, solution.RetryNum, solution.TopicSolutions)
	doc, err := do.toDoc()
	if err != nil {
		return err
	}

	filter, err := impl.dao.DocIdFilter(solution.Id)
	if err != nil {
		return err
	}

	_, err = impl.dao.ReplaceDoc(filter, doc)

	return err
}
