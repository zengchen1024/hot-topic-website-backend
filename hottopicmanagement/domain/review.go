package domain

import (
	"errors"

	"github.com/opensourceways/hot-topic-website-backend/utils"
)

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

func (r *TopicToReview) NewHotTopic(date string) HotTopic {
	dss := make([]DiscussionSource, len(r.DiscussionSources))
	for i := range r.DiscussionSources {
		item := &r.DiscussionSources[i].DiscussionSource
		item.ImportedAt = date

		dss[i] = *item
	}

	logTime := date
	if t, err := findMaxDate(dss); err == nil {
		logTime = utils.GetDate(&t)
	}

	return HotTopic{
		Title: r.Title,
		TransferLogs: []TransferLog{
			TransferLog{
				Order: r.Order,
				Date:  date,
				StatusLog: StatusLog{
					Time:   logTime,
					Status: statusNew,
				},
			},
		},
		DiscussionSources: dss,
	}
}

func (r *TopicToReview) getAppendedDS() []DiscussionSource {
	v := make([]DiscussionSource, 0, len(r.DiscussionSources))

	for i := range r.DiscussionSources {
		if item := &r.DiscussionSources[i]; !item.isOldOne() {
			v = append(v, item.DiscussionSource)
		}
	}

	return v
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

func NewTopicsToReview() TopicsToReview {
	return TopicsToReview{
		Candidates: map[string][]TopicToReview{},
	}
}

func (t *TopicsToReview) CandidatesNum() int {
	n := 0

	for i := range t.Candidates {
		n += len(t.Candidates[i])
	}

	return n
}

func (t *TopicsToReview) UpdateSelected(lastHotTopic string, items []TopicToReview) error {
	lastTopics := map[string]*TopicToReview{}
	for i := range t.Selected {
		if item := &t.Selected[i]; item.Category == lastHotTopic {
			lastTopics[item.Title] = item
		}
	}

	n := 0
	for i := range items {
		if old, ok := lastTopics[items[i].Title]; ok {
			n++

			if err := old.checkIfOldDSMissing(&items[i]); err != nil {
				return err
			}
		}
	}

	if n != len(lastTopics) {
		return errors.New("missing last hot topics")
	}

	t.Selected = items

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
