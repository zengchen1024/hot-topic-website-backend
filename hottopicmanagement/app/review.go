package app

import (
	"fmt"

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

func (s *appService) ApplyToHotTopic(community string, date string) error {
	review, err := s.repoTopicsToReview.FindSelected(community)
	if err != nil {
		return err
	}

	selectedMap := make(map[string]*domain.TopicToReview, len(review.Selected))
	for i := range review.Selected {
		item := &review.Selected[i]
		selectedMap[item.Title] = item
	}

	hts, err := s.repoHotTopic.FindOpenOnes(community)
	if err != nil {
		return err
	}

	for i := range hts {
		old := &hts[i]

		r := selectedMap[old.Title]
		if r == nil {
			return fmt.Errorf("no corresponding topic to review for the hot topic(%s)", old.Id)
		}

		old.Update(r, date)

		delete(selectedMap, old.Title)
	}

	newOnes := make([]domain.HotTopic, len(selectedMap))
	if len(selectedMap) > 0 {
		i := 0
		for _, r := range selectedMap {
			newOnes[i] = r.NewHotTopic(date)
			i++
		}
	}

	for i := range hts {
		if err := s.repoHotTopic.Save(community, &hts[i]); err != nil {
			return err
		}
	}

	for i := range newOnes {
		if err := s.repoHotTopic.Add(community, &newOnes[i]); err != nil {
			return err
		}
	}

	return nil
}
