package domain

import (
	"fmt"

	"github.com/opensourceways/hot-topic-website-backend/utils"
)

type DiscussionSource struct {
	Id         int    `json:"id"           required:"true"`
	URL        string `json:"url"          required:"true"`
	Type       string `json:"source_type"  required:"true"`
	SourceId   string `json:"source_id"    required:"true"`
	CreatedAt  string `json:"created_at"   required:"true"`
	ImportDate string
}

func (ds *DiscussionSource) isOldOne() bool {
	return ds.ImportDate != ""
}

type StatusLog struct {
	Status string
	Time   string
}

func (s *StatusLog) resolved() bool {
	return s.Status == "Resolved"
}

type TransferLog struct {
	StatusLog
	Order int    // the topic is ordered on the report of that week
	Date  string // the date that the report of that week is created
}

// HotTopic
type HotTopic struct {
	Id                string
	Title             string
	Order             int
	TransferLogs      []TransferLog
	DiscussionSources []DiscussionSource
	StatusTransferLog []StatusLog //TODO: delete
	Version           int
}

func (ht *HotTopic) CheckIfIsAGoodReview(t *TopicToReview) error {
	m := t.GetDSSet()

	for i := range ht.DiscussionSources {
		if v := ht.DiscussionSources[i].Id; !m[v] {
			return fmt.Errorf(
				"missing discussion source(%d) for the reviewing topic(%s)", v, t.Title,
			)
		}
	}

	return nil
}

func NewHotTopic(title string, order int, sources []DiscussionSource, createdAt string) HotTopic {
	return HotTopic{
		Title:             title,
		Order:             order,
		DiscussionSources: sources,
		TransferLogs: []TransferLog{
			{
				StatusLog: StatusLog{
					Time:   createdAt,
					Status: "New",
				},
				Order: order,
				Date:  utils.Date(),
			},
		},
	}
}

func (ht *HotTopic) GetDSSet() map[int]bool {
	v := make(map[int]bool, len(ht.DiscussionSources))

	for i := range ht.DiscussionSources {
		v[ht.DiscussionSources[i].Id] = true
	}

	return v
}

func (ht *HotTopic) IsResolved() bool {
	for i := len(ht.StatusTransferLog) - 1; i >= 0; i-- {
		if ht.StatusTransferLog[i].resolved() {
			return true
		}
	}

	return false
}

func (ht *HotTopic) GetDiscussionSource(dsId int) *DiscussionSource {
	for i := range ht.DiscussionSources {
		if item := &ht.DiscussionSources[i]; item.Id == dsId {
			return item
		}
	}

	return nil
}

// DiscussionSourceInfo
type DiscussionSourceInfo struct {
	Id    int
	URL   string
	Title string

	removed bool
}

func (info *DiscussionSourceInfo) Removed() bool {
	return info.removed
}

// NotHotTopic
type NotHotTopic struct {
	Title             string
	DiscussionSources []DiscussionSourceInfo
}

func NewNotHotTopic(title string, sources []DiscussionSourceInfo) NotHotTopic {
	return NotHotTopic{
		Title:             title,
		DiscussionSources: sources,
	}
}

func (nht *NotHotTopic) GetDSSet() map[int]bool {
	v := make(map[int]bool, len(nht.DiscussionSources))

	for i := range nht.DiscussionSources {
		v[nht.DiscussionSources[i].Id] = true
	}

	return v
}

func (nht *NotHotTopic) UpdateRemoved(dsIdsOfNewTopic map[int]bool) {
	for i := range nht.DiscussionSources {
		item := &nht.DiscussionSources[i]

		if _, ok := dsIdsOfNewTopic[item.Id]; !ok {
			item.removed = true
		}
	}
}

func (nht *NotHotTopic) Sort() []*DiscussionSourceInfo {
	v := make([]*DiscussionSourceInfo, len(nht.DiscussionSources))
	h := 0
	t := len(v) - 1
	for i := range nht.DiscussionSources {
		if item := &nht.DiscussionSources[i]; item.removed {
			v[h] = item
			h++
		} else {
			v[t] = item
			t--
		}
	}

	return v
}
