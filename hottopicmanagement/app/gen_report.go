package app

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

const (
	keyHotTopics = "hot_topics"
)

func (s *appService) GenReport(community string) (string, error) {
	dbReport, err := s.repoTopicReport.GetCurrentReport(community)
	if err != nil {
		return "Get DB report failed.", err
	}
	var report map[string]interface{}
	mapstructure.Decode(dbReport, &report)
	topics := []domain.HotTopic{}
	for _, top := range dbReport.TopN {
		topic, err := s.repoHotTopic.FindOneWithId(community, top.TopicId)
		if err != nil {
			continue
		}
		topics = append(topics, topic)
	}
	report[keyHotTopics] = topics

	jsonData, err := json.Marshal(report)
	if err != nil {
		return "Data to json failed.", err
	}

	// return json string
	return string(jsonData), nil
}
