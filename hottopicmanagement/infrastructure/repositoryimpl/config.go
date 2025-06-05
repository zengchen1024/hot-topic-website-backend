package repositoryimpl

type Config struct {
	Collections Collections `json:"collections"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Collections,
	}
}

type Collections struct {
	HotTopic    string `json:"hot_topic"       required:"true"`
	NotHotTopic string `json:"not_hot_topic"   required:"true"`
}
