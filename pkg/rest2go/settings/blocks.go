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

type Server struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	HealthCheck     bool   `yaml:"health-check"`
	NotFoundHandler bool   `yaml:"not-found-handler"`
}

func newServer() Server {
	return Server{
		Host:            "0.0.0.0",
		Port:            8080,
		HealthCheck:     false,
		NotFoundHandler: false,
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

type Database struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Schema   string `yaml:"schema"`
}

func newDatabase() Database {
	return Database{
		Driver:   "sqlite3",
		Host:     "./data/app.db",
		Port:     0,
		User:     "",
		Password: "",
		Name:     "",
		Schema:   "",
	}
}
