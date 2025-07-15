package app

type Config struct {
	FilePath   string `json:"filt_path"      required:"true"`
	SaveToFile bool   `json:"save_to_file"`
}

func (cfg *Config) SetDefault() {
	cfg.SaveToFile = true
}
