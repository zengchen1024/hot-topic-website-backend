package repositoryimpl

import (
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

func (impl *hotTopic) Resolved(community string, topic *domain.HotTopic, closedAt int64) error {
	do := tohotTopicDO(topic)
	doc, err := do.toDoc()
	if err != nil {
		return err
	}
	doc[fieldClosedAt] = closedAt

	dao, err := impl.dao(community)
	if err != nil {
		return err
	}

	docFilter, err := dao.DocIdFilter(topic.Id)
	if err != nil {
		return err
	}

	return dao.UpdateDocsWithoutVersion(docFilter, doc)
}

func (impl *hotTopic) Appended(community string, topic *domain.HotTopic) error {
	do := tohotTopicDO(topic)
	doc, err := do.toDoc()
	if err != nil {
		return err
	}

	dao, err := impl.dao(community)
	if err != nil {
		return err
	}

	docFilter, err := dao.DocIdFilter(topic.Id)
	if err != nil {
		return err
	}

	return dao.UpdateDocsWithoutVersion(docFilter, doc)
}
