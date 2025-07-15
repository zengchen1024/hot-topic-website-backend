package app

import (
	"fmt"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
	"github.com/opensourceways/hot-topic-website-backend/utils"
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

	if err := t.UpdateSelected(cmd.Selected); err != nil {
		return err
	}

	// TODO check if all the discussion sources are valid, for example there are from Data Clean

	return s.repoTopicsToReview.SaveSelected(community, &t)
}

func (s *appService) GetHotTopics(community string, date int64) (HotTopicsDTO, error) {
	hts, err := s.repoHotTopic.FindAll(community, date)
	if err != nil {
		return HotTopicsDTO{}, err
	}

	return toHotTopicsDTO(hts, date), nil
}

func (s *appService) getReviews(community string, dateSec int64) (review domain.TopicsToReview, err error) {
	review, err = s.repoTopicsToReview.Find(community)
	if err != nil {
		return
	}

	if !review.IsMatchedReview(dateSec) {
		err = fmt.Errorf(
			"review is not right one which match the time. expect:%d, has:%d",
			dateSec, review.Date,
		)
	}

	return
}

func (s *appService) GetTopicsToPublish(community string) (dto HotTopicsDTO, err error) {
	date := utils.GetLastFriday()
	dateSec := date.Unix()

	review, err := s.getReviews(community, dateSec)
	if err != nil {
		return
	}

	hts, err := s.repoHotTopic.FindAll(community, dateSec)
	if err != nil {
		return
	}

	_, news, err := review.FilterChangedAndNews(hts, date)
	if err != nil {
		return
	}

	dto = toHotTopicsDTO(append(hts, news...), dateSec)

	return
}

func (s *appService) ApplyToHotTopic(community string) error {
	date := utils.GetLastFriday()
	dateSec := date.Unix()

	review, err := s.getReviews(community, dateSec)
	if err != nil {
		return err
	}

	hts, err := s.repoHotTopic.FindAll(community, dateSec)
	if err != nil {
		return err
	}

	changed, news, err := review.FilterChangedAndNews(hts, date)
	if err != nil {
		return err
	}

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
