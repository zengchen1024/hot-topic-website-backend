package repositoryimpl

type Config struct {
	TopicReview          string                 `json:"topic_review"           required:"true"`
	TopicSolution        string                 `json:"topic_solution"         required:"true"`
	CommunityCollections []CommunityCollections `json:"Community_collections"`
}

func (cfg *Config) ConfigItems() []interface{} {
	r := make([]interface{}, len(cfg.CommunityCollections))

	for i := range cfg.CommunityCollections {
		r[i] = &cfg.CommunityCollections[i]
	}

	return r
}

type CommunityCollections struct {
	Community   string      `json:"community"    required:"true"`
	Collections Collections `json:"collections"`
}

func (cfg *CommunityCollections) ConfigItems() []interface{} {
	return []interface{}{&cfg.Collections}
}

type Collections struct {
	HotTopic    string `json:"hot_topic"       required:"true"`
	NotHotTopic string `json:"not_hot_topic"   required:"true"`
}
