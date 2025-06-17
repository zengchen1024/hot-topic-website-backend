package app

import (
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

type CmdToUploadOptionalTopics []OptionalTopic

func (cmd CmdToUploadOptionalTopics) init() {
	for i := range cmd {
		cmd[i].init()
	}
}

type DiscussionSourceInfos []DiscussionSourceInfo

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

// DiscussionSourceInfo
type DiscussionSourceInfo struct {
	Title  string `json:"title"          required:"true"`
	Closed bool   `json:"source_closed"  required:"true"`

	domain.DiscussionSource

	appended bool // if true, it is newly appended to the old hot topic
}
