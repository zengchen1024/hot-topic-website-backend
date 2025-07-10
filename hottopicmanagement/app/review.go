package app

import (
	"time"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

func (s *appService) toSelected(
	oldTopics []domain.HotTopic, current map[string]*OptionalTopic,
) (
	[]domain.TopicToReview, error,
) {
	r := make([]domain.TopicToReview, len(oldTopics))

	for i := range oldTopics {
		old := &oldTopics[i]

		item := current[old.Title]
		item.updateAppended(old.GetDSSet())

		r[i] = item.toTopicToReview()

		if err := old.InitReview(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (s *appService) GetTopicsToReview(community string) (TopicsToReviewDTO, error) {
	return s.repoTopicsToReview.Find(community)
}

func (s *appService) UpdateSelected(community string, cmd *CmdToUpdateSelected) error {
	t, err := s.repoTopicsToReview.FindSelected(community)
	if err != nil {
		return err
	}

	if err := t.UpdateSelected(sheetLastTopics, cmd.Selected); err != nil {
		return err
	}

	// TODO check if all the discussion sources are valid, for example there are from Data Clean

	return s.repoTopicsToReview.SaveSelected(community, &t)
}

func (s *appService) GetTopicsToPublish(community string) (TopicsToPublishDTO, error) {
	v, err := s.repoTopicsToReview.FindSelected(community)
	if err != nil {
		return TopicsToPublishDTO{}, err
	}

	return TopicsToPublishDTO{v.Selected}, nil
}

func (s *appService) GetHotTopics(community string, since int64) (HotTopicsDTO, error) {
	hts, err := s.repoHotTopic.FindAll(community, since)
	if err != nil {
		return HotTopicsDTO{}, err
	}

	items := make([]hotTopicDTO, len(hts))
	for i := range hts {
		item := &hts[i]

		items[i] = hotTopicDTO{
			Title:             item.Title,
			DiscussionSources: item.DiscussionSources,
		}
	}

	return HotTopicsDTO{Topics: items}, nil
}

func (s *appService) ApplyToHotTopic(community string, date time.Time) error {
	review, err := s.repoTopicsToReview.Find(community)
	if err != nil {
		return err
	}

	since := date.AddDate(0, -1, 0)
	hts, err := s.repoHotTopic.FindAll(community, since.Unix())
	if err != nil {
		return err
	}

	changed, news := review.FilterChangedAndNews(hts, date)

	for i := range changed {
		if err := s.repoHotTopic.Save(community, changed[i]); err != nil {
			return err
		}
	}

	for i := range news {
		if err := s.repoHotTopic.Add(community, &news[i]); err != nil {
			return err
		}
	}

	return s.repoNotHotTopic.Save(community, review.GenNotHotTopics())
}
