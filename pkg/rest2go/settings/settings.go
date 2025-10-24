package settings

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

var (
	loadOnce sync.Once
	cached   any
	loadErr  error
	paths    = [2]string{"./app.yml", "./configs/app.yml"}
)

type Defaults interface {
	SetDefaults()
	ConfigureLogs()
}

type Settings struct {
	Logs          Logs          `yaml:"logs"`
	Server        Server        `yaml:"server"`
	Authorization Authorization `yaml:"authorization"`
}

func (s *Settings) SetDefaults() {
	s.Logs = newLogs()
	s.Server = newServer()
	s.Authorization = newAuthorization()
}

func (s *Settings) ConfigureLogs() {
	var handler *slog.TextHandler
	var level slog.Level

	switch s.Logs.Level {
	case "debug":
		level = slog.LevelDebug

	case "warn":
		level = slog.LevelWarn

	case "error":
		level = slog.LevelError

	default:
		level = slog.LevelInfo
	}

	options := &slog.HandlerOptions{
		Level: level,
	}

	if s.Logs.FileEnabled {
		file := &lumberjack.Logger{
			Filename:  "./logs/app.log",
			MaxSize:   s.Logs.MaxSize,
			MaxAge:    s.Logs.MaxAge,
			LocalTime: true,
		}
		writer := io.MultiWriter(os.Stdout, file)
		handler = slog.NewTextHandler(writer, options)
	} else {
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func Load[T any]() (T, error) {
	var empty T

	loadOnce.Do(func() {
		var settings T
		target, cast := any(&settings).(Defaults)

		if !cast {
			loadErr = fmt.Errorf("defaults not implemented for application settings structure")
			return
		}

		target.SetDefaults()
		var file *os.File
		var fileErr error

		for _, path := range paths {
			file, fileErr = os.Open(path)

			if fileErr != nil {
				slog.Warn("Settings not found", "path", path)
				continue
			}

			slog.Info("Settings found", "path", path)
			break
		}

		if file != nil {
			defer file.Close()
			decoder := yaml.NewDecoder(file)

			if err := decoder.Decode(&settings); err != nil {
				loadErr = err
				return
			}
		}

		target.ConfigureLogs()
		cached = settings
	})

	if loadErr != nil {
		return empty, loadErr
	}

	if cfg, cast := cached.(T); cast {
		return cfg, nil
	}

	return empty, fmt.Errorf("application settings already loaded as different type")
}
