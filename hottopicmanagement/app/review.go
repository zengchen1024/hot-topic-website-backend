package app

import (
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

		if err := old.CheckIfIsAGoodReview(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (s *appService) GetTopicsToReview(community string) (TopicsToReviewDTO, error) {
	return s.repoTopicsToReview.Find(community)
}

func (s *appService) UpdateSelected(cmd *CmdToUpdateSelected) error {
	t, err := s.repoTopicsToReview.FindSelected(cmd.Community)
	if err != nil {
		return err
	}

	if err := t.UpdateSelected(sheetLastTopics, cmd.Selected); err != nil {
		return err
	}

	// TODO check if all the discussion sources are valid, for example there are from Data Clean

	return s.repoTopicsToReview.SaveSelected(cmd.Community, &t)
}

func (s *appService) GetTopicsToPublish(community string) ([]domain.TopicToReview, error) {
	v, err := s.repoTopicsToReview.FindSelected(community)
	if err != nil {
		return nil, err
	}

	return v.Selected, nil
}
