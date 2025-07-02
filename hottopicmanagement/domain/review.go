package domain

import "errors"

// DiscussionSourceToReview
type DiscussionSourceToReview struct {
	Title  string `json:"title"`
	Closed bool   `json:"source_closed"`

	DiscussionSource
}

// TopicToReview
type TopicToReview struct {
	Order             int                        `json:"order"`
	Title             string                     `json:"title"`
	Category          string                     `json:"category"`
	Resolved          bool                       `json:"resolved"`
	DiscussionSources []DiscussionSourceToReview `json:"dss"`
}

func (t *TopicToReview) GetDSSet() map[int]bool {
	v := make(map[int]bool, len(t.DiscussionSources))

	for i := range t.DiscussionSources {
		v[t.DiscussionSources[i].Id] = true
	}

	return v
}

func (t *TopicToReview) getDSMap() map[int]*DiscussionSourceToReview {
	v := make(map[int]*DiscussionSourceToReview, len(t.DiscussionSources))

	for i := range t.DiscussionSources {
		item := &t.DiscussionSources[i]
		v[item.Id] = item
	}

	return v
}

func (t *TopicToReview) getOldDS() map[int]bool {
	oldOnes := make(map[int]bool, len(t.DiscussionSources))

	for i := range t.DiscussionSources {
		if item := &t.DiscussionSources[i]; item.isOldOne() {
			oldOnes[item.Id] = true
		}
	}

	return oldOnes
}

func (t *TopicToReview) checkIfOldDSMissing(t1 *TopicToReview) error {
	oldOne := t.getOldDS()
	oldOne1 := t1.getOldDS()

	err := errors.New("missing old ds")

	if len(oldOne) != len(oldOne1) {
		return err
	}

	for k := range oldOne {
		if !oldOne1[k] {
			return err
		}
	}

	return nil
}

// TopicsToReview
type TopicsToReview struct {
	Candidates map[string][]TopicToReview `json:"cadidates"`
	Selected   []TopicToReview            `json:"selected"`
	Version    int                        `json:"-"`
}

func (t *TopicsToReview) CandidatesNum() int {
	n := 0

	for i := range t.Candidates {
		n += len(t.Candidates[i])
	}

	return n
}

func (t *TopicsToReview) UpdateSelected(lastHotTopic string, items []TopicToReview) error {
	m := map[string]*TopicToReview{}
	for i := range t.Selected {
		if item := &t.Selected[i]; item.Category == lastHotTopic {
			m[item.Title] = item
		}
	}

	n := 0
	for i := range items {
		if old, ok := m[items[i].Title]; ok {
			n++

			if err := old.checkIfOldDSMissing(&items[i]); err != nil {
				return err
			}
		}
	}

	if n != len(m) {
		return errors.New("missing last hot topics")
	}

	return nil
}

func (t *TopicsToReview) SetSelected(cantegory string, items []TopicToReview) {
	for i := range items {
		items[i].Category = cantegory
	}

	t.Selected = items
}

func (t *TopicsToReview) AddCandidate(category string, topic *TopicToReview) {
	topic.Category = category

	t.Candidates[category] = append(t.Candidates[category], *topic)
}
