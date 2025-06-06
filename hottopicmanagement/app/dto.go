package app

import (
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

type CmdToUploadOptionalTopics []OptionalTopic

// OptionalTopic
type OptionalTopic struct {
	Title             string                 `json:"summary"    required:"true"`
	DiscussionSources []DiscussionSourceInfo `json:"discussion" required:"true"`
}

func (ot *OptionalTopic) updateAppended(dsIdsOfOldTopic map[int]bool) {
	for i := range ot.DiscussionSources {
		item := &ot.DiscussionSources[i]

		if _, ok := dsIdsOfOldTopic[item.Id]; !ok {
			item.appended = true
		}
	}
}

func (ot *OptionalTopic) sort() []*DiscussionSourceInfo {
	v := make([]*DiscussionSourceInfo, len(ot.DiscussionSources))
	h := 0
	t := len(v) - 1
	for i := range ot.DiscussionSources {
		if item := &ot.DiscussionSources[i]; item.appended {
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
	v := make(map[int]bool, len(ot.DiscussionSources))

	for i := range ot.DiscussionSources {
		v[ot.DiscussionSources[i].Id] = true
	}

	return v
}

// DiscussionSourceInfo
type DiscussionSourceInfo struct {
	Title string `json:"title" required:"true"`
	domain.DiscussionSource

	appended bool // if true, it is newly appended to the old hot topic
}
