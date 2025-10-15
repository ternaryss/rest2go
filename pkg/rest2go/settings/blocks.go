package settings

type Logs struct {
	Level       string `yaml:"level"`
	FileEnabled bool   `yaml:"file-enabled"`
	MaxSize     int    `yaml:"max-size"`
	MaxAge      int    `yaml:"max-age"`
}

func newLogs() Logs {
	return Logs{
		Level:       "info",
		FileEnabled: false,
		MaxSize:     10,
		MaxAge:      7,
	}
}
