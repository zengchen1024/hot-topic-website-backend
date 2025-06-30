package domain

type DiscussionSource struct {
	Id        int    `json:"id"           required:"true"`
	URL       string `json:"url"          required:"true"`
	Type      string `json:"source_type"  required:"true"`
	SourceId  string `json:"source_id"    required:"true"`
	CreatedAt string `json:"created_at"   required:"true"`
}

type StatusLog struct {
	Status string
	Time   string
}

// HotTopic
type HotTopic struct {
	Id                string             `json:id`
	Title             string             `json:title`
	Order             int                `json:order`
	DiscussionSources []DiscussionSource `json:discussion`
	StatusTransferLog []StatusLog        `json:transferlog`
	Version           int                `json:version`
}

// TopicReport
type TopicTopN struct {
	Idx     int    `json:"idx" bson:"idx"`
	TopicId string `json:"topic_id" bson:"topic"`
}
type TopicReport struct {
	Year int         `json:"year" bson:"year"`
	Week int         `json:"week" bson:"week"`
	Cnt  int         `json:"cnt" bson:"cnt"`
	TopN []TopicTopN `json:"topn" bson:"topn"`
}

func NewHotTopic(title string, order int, sources []DiscussionSource, createdAt string) HotTopic {
	return HotTopic{
		Title:             title,
		Order:             order,
		DiscussionSources: sources,
		StatusTransferLog: []StatusLog{
			{
				Time:   createdAt,
				Status: "New",
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
