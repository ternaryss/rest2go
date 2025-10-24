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

type Header struct {
	Enabled bool     `yaml:"enabled"`
	Key     string   `yaml:"key"`
	Public  []string `yaml:"public"`
}

func newHeader() Header {
	return Header{
		Enabled: false,
		Key:     "",
		Public:  []string{},
	}
}

type Authorization struct {
	Header Header `yaml:"header"`
}

func newAuthorization() Authorization {
	return Authorization{
		Header: newHeader(),
	}
}
