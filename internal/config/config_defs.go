package config

type Config struct {
	Main          `mapstructure:"main"`
	Proxy         Proxy         `mapstructure:"proxy"`
	RenameRule    RenameRule    `mapstructure:"rename_rule"`
	LoggerOptions LoggerOptions `mapstructure:"logger_options"`
}

type Main struct {
	Mode                   string `mapstructure:"mode"`
	SourceDirectory        string `mapstructure:"source_directory"`
	SuccessOutputDirectory string `mapstructure:"success_output_directory"`
}

type PosterWaterMark struct {
	Enabled  bool   `mapstructure:"enabled"`
	Position string `mapstructure:"position"`
}

type Proxy struct {
	Enabled bool   `mapstructure:"enabled"`
	URL     string `mapstructure:"url"`
	Timeout int    `mapstructure:"timeout"`
	Retry   int    `mapstructure:"retry"`
}
type RenameRule struct {
	LocationRule        string `mapstructure:"location_rule"`
	FileRule            string `mapstructure:"file_rule"`
	RenameImgWithNumber bool   `mapstructure:"rename_img_with_number"`
}

type LoggerOptions struct {
	Level   string `mapstructure:"level"`
	LogPath string `mapstructure:"log_path"`
}
