package settings

type Logs struct {
	Level       string `yaml:"level" json:"level"`
	FileEnabled bool   `yaml:"file-enabled" json:"fileEnabled"`
	MaxSize     int    `yaml:"max-size" json:"maxSize"`
	MaxAge      int    `yaml:"max-age" json:"maxAge"`
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
	Host            string `yaml:"host" json:"host"`
	Port            int    `yaml:"port" json:"port"`
	Configuration   bool   `yaml:"configuration" json:"configuration"`
	HealthCheck     bool   `yaml:"health-check" json:"healthCheck"`
	NotFoundHandler bool   `yaml:"not-found-handler" json:"notFoundHandler"`
}

func newServer() Server {
	return Server{
		Host:            "0.0.0.0",
		Port:            8080,
		Configuration:   false,
		HealthCheck:     false,
		NotFoundHandler: false,
	}
}

type Header struct {
	Enabled bool     `yaml:"enabled" json:"enabled"`
	Key     string   `yaml:"key" json:"-"`
	Public  []string `yaml:"public" json:"public"`
}

func newHeader() Header {
	return Header{
		Enabled: false,
		Key:     "",
		Public:  []string{},
	}
}

type Authorization struct {
	Header Header `yaml:"header" json:"header"`
}

func newAuthorization() Authorization {
	return Authorization{
		Header: newHeader(),
	}
}

type Database struct {
	Driver   string `yaml:"driver" json:"driver"`
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"-"`
	Name     string `yaml:"name" json:"name"`
	Schema   string `yaml:"schema" json:"schema"`
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
