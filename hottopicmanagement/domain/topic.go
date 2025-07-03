package domain

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/hot-topic-website-backend/utils"
)

const (
	statusNew      = "New"
	statusAppended = "Appended"
	statusResolved = "Resolved"
)

type DiscussionSourceMeta struct {
	Id        int    `json:"id"           required:"true"`
	URL       string `json:"url"          required:"true"`
	Type      string `json:"source_type"  required:"true"`
	SourceId  string `json:"source_id"    required:"true"`
	CreatedAt string `json:"created_at"   required:"true"`
}

type DiscussionSource struct {
	DiscussionSourceMeta

	ImportedAt string `json:"imported_at"`
}

func (ds *DiscussionSource) isOldOne() bool {
	return ds.ImportedAt != ""
}

func (ds *DiscussionSource) setImportDate(date string) {
	if !ds.isOldOne() {
		ds.ImportedAt = date
	}
}

func findMaxDate(items []DiscussionSource) (time.Time, error) {
	if len(items) == 0 {
		return time.Time{}, fmt.Errorf("no dss")
	}

	maxTime, err := time.Parse(time.RFC3339, items[0].CreatedAt)
	if err != nil {
		return time.Time{}, err
	}

	items = items[1:]
	for i := range items {
		t, err := time.Parse(time.RFC3339, items[i].CreatedAt)
		if err != nil {
			return time.Time{}, err
		}

		if t.After(maxTime) {
			maxTime = t
		}
	}

	return maxTime, nil
}

type StatusLog struct {
	Status string
	Time   string
}

func (s *StatusLog) resolved() bool {
	return s.Status == statusResolved
}

type TransferLog struct {
	StatusLog
	Order int    // the topic is ordered on the report of that week
	Date  string // the date that the report of that week is created
}

func newAppendedLog(items []DiscussionSource, date string, aWeekAgo time.Time) (log StatusLog) {
	log.Status = statusAppended

	if v, err := findMaxDate(items); err == nil && v.After(aWeekAgo) {
		log.Time = utils.GetDate(&v)
	} else {
		logrus.Errorf("find max date failed, err:%s", err.Error())

		log.Time = date
	}

	return
}

// HotTopic
type HotTopic struct {
	Id                string
	Title             string
	TransferLogs      []TransferLog
	DiscussionSources []DiscussionSource
	Version           int
}

func (ht *HotTopic) Order() int {
	if n := len(ht.TransferLogs); n > 0 {
		return ht.TransferLogs[n-1].Order
	}

	return 0
}

func (ht *HotTopic) Update(r *TopicToReview, date string, aWeekAgo time.Time) {
	logNum := len(ht.TransferLogs)
	if logNum == 0 {
		// it is impossible that there aren't old logs
		return
	}

	if ht.TransferLogs[logNum-1].Date == date {
		logrus.Info("it is repeated to update the hot topic")

		// it is repeated to update the hot topic
		return
	}

	items := r.getAppendedDS()
	if len(items) > 0 {
		for i := range items {
			items[i].ImportedAt = date
		}

		ht.DiscussionSources = append(ht.DiscussionSources, items...)
	}

	log := TransferLog{
		Date:  date,
		Order: r.Order,
	}
	if r.Resolved {
		log.Status = statusResolved
		log.Time = date
	} else if len(items) > 0 {
		log.StatusLog = newAppendedLog(items, date, aWeekAgo)
	} else {
		log.StatusLog = ht.TransferLogs[len(ht.TransferLogs)-1].StatusLog
	}

	ht.TransferLogs = append(ht.TransferLogs, log)
}

func (ht *HotTopic) InitReview(t *TopicToReview) error {
	m := t.getDSMap()

	for i := range ht.DiscussionSources {
		item := ht.DiscussionSources[i]

		v, ok := m[item.Id]
		if !ok {
			return fmt.Errorf(
				"missing discussion source(%d) for the reviewing topic(%s)", v, t.Title,
			)
		}
		v.ImportedAt = item.ImportedAt
	}

	t.Order = ht.Order()

	return nil
}

func (ht *HotTopic) GetDSSet() map[int]bool {
	v := make(map[int]bool, len(ht.DiscussionSources))

	for i := range ht.DiscussionSources {
		v[ht.DiscussionSources[i].Id] = true
	}

	return v
}

func (ht *HotTopic) IsResolved() bool {
	n := len(ht.TransferLogs)

	return n > 0 && ht.TransferLogs[n-1].resolved()
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
