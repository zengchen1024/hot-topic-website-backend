package app

import (
	"errors"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

// CmdToUploadOptionalTopics
type CmdToUploadOptionalTopics []OptionalTopic

func (cmd CmdToUploadOptionalTopics) init() {
	for i := range cmd {
		cmd[i].init()
	}
}

// OptionalTopic
type OptionalTopic struct {
	Title             string                  `json:"summary"    required:"true"`
	DiscussionSources []DiscussionSourceInfos `json:"discussion" required:"true"`

	discussionSources []*DiscussionSourceInfo
	total             int
}

func (ot *OptionalTopic) init() {
	n := 0
	for i := range ot.DiscussionSources {
		n += len(ot.DiscussionSources[i])
	}
	ot.total = n

	items := make([]*DiscussionSourceInfo, n)
	k := 0
	for i := range ot.DiscussionSources {
		s1 := ot.DiscussionSources[i]
		for j := range s1 {
			items[k] = &s1[j]
			k++
		}
	}
	ot.discussionSources = items
}

func (ot *OptionalTopic) updateAppended(dsIdsOfOldTopic map[int]bool) {
	for i := range ot.discussionSources {
		item := ot.discussionSources[i]

		if _, ok := dsIdsOfOldTopic[item.Id]; !ok {
			item.appended = true
		}
	}
}

func (ot *OptionalTopic) sort() []*DiscussionSourceInfo {
	v := make([]*DiscussionSourceInfo, ot.total)
	h := 0
	t := len(v) - 1
	for i := range ot.discussionSources {
		if item := ot.discussionSources[i]; item.appended {
			v[h] = item
			h++
		} else {
			v[t] = item
			t--
		}
	}

	return v
}

func (ot *OptionalTopic) getDSSet() map[int]bool {
	v := make(map[int]bool, ot.total)

	for i := range ot.discussionSources {
		v[ot.discussionSources[i].Id] = true
	}

	return v
}

func (ot *OptionalTopic) toTopicToReview() (t domain.TopicToReview) {
	t.Title = ot.Title

	v := make([]domain.DiscussionSourceToReview, ot.total)
	items := ot.sort()
	for i := range items {
		v[i] = items[i].toDiscussionSourceToReview()
	}

	t.DiscussionSources = v

	return
}

// DiscussionSourceInfos
type DiscussionSourceInfos []DiscussionSourceInfo

func (infos DiscussionSourceInfos) filterout() (resolved, unresolved []*DiscussionSourceInfo) {
	for i := range infos {
		if item := &infos[i]; item.Closed {
			resolved = append(resolved, item)
		} else {
			unresolved = append(unresolved, item)
		}
	}

	return
}

func (infos DiscussionSourceInfos) resolvedNum() int {
	num := 0
	for i := range infos {
		if infos[i].Closed {
			num++
		}
	}

	return num
}

func (infos DiscussionSourceInfos) sort() []*DiscussionSourceInfo {
	v := make([]*DiscussionSourceInfo, len(infos))
	h := 0
	t := len(v) - 1
	for i := range infos {
		if item := &infos[i]; item.appended {
			v[h] = item
			h++
		} else {
			v[t] = item
			t--
		}
	}

	return v
}

// DiscussionSourceInfo
type DiscussionSourceInfo struct {
	Closed bool `json:"source_closed"  required:"true"`

	domain.DiscussionSourceMeta

	// if true, it is newly appended to the old hot topic
	appended bool
}

func (info *DiscussionSourceInfo) toDiscussionSourceToReview() domain.DiscussionSourceToReview {
	return domain.DiscussionSourceToReview{
		Closed: info.Closed,
		DiscussionSource: domain.DiscussionSource{
			DiscussionSourceMeta: info.DiscussionSourceMeta,
		},
	}
}

// TopicsToReviewDTO
type TopicsToReviewDTO = domain.TopicsToReview

// CmdToUpdateSelected
type CmdToUpdateSelected struct {
	Selected []domain.TopicToReview `json:"selected"`
}

func (cmd *CmdToUpdateSelected) Validate() error {
	if err := cmd.checkDuplicateDS(); err != nil {
		return err
	}

	return cmd.checkDuplicateTopic()
}

func (cmd *CmdToUpdateSelected) checkDuplicateTopic() error {
	m := make(map[string]bool, len(cmd.Selected))
	for i := range cmd.Selected {
		t := cmd.Selected[i].Title
		if m[t] {
			return errors.New("there are duplicate topics")
		}
		m[t] = true
	}

	return nil
}

func (cmd *CmdToUpdateSelected) checkDuplicateDS() error {
	m := make(map[int]bool, len(cmd.Selected))
	for i := range cmd.Selected {
		items := cmd.Selected[i].DiscussionSources
		for j := range items {
			v := items[j].Id
			if m[v] {
				return errors.New("there are duplicate discussion sources")
			}
			m[v] = true
		}
	}

	return nil
}

type TopicsToPublishDTO struct {
	Topics []domain.TopicToReview `json:"topics"`
}

type HotTopicsDTO struct {
	Topics []hotTopicDTO `json:"topics"`
}

type hotTopicDTO struct {
	Id                string                    `json:"id"`
	Order             int                       `json:"order"`
	Title             string                    `json:"title"`
	Status            statusLogDTO              `json:"status"`
	DiscussionSources []domain.DiscussionSource `json:"dss"`
}

type statusLogDTO struct {
	Time   string `json:"time"`
	Status string `json:"status"`
}
