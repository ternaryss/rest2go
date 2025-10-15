package file

import (
	"testing"

	"github.com/ternaryss/rest2go/pkg/rest2go/settings"
)

type AppSettings struct {
	settings.Settings `yaml:",inline"`
	Foo               string `yaml:"foo"`
}

func (s *AppSettings) SetDefaults() {
	s.Settings.SetDefaults()
	s.Foo = "bar"
}

func TestLoad(tst *testing.T) {
	// Given
	foo := "testing"
	logsLevel := "debug"
	logsFileEnabled := true
	logsMaxSize := 20
	logsMaxAge := 30

	// When
	settings, err := settings.Load[AppSettings]()

	// Then
	if err != nil {
		tst.Errorf("Loading settings failed: %s", err)
	}

	logs := settings.Logs

	if settings.Foo != foo {
		tst.Errorf("Invalid foo - expected: %s, value: %s", foo, settings.Foo)
	}

	if logs.Level != logsLevel {
		tst.Errorf("Invalid logs.level - expected: %s, value: %s", logsLevel, logs.Level)
	}

	if logs.FileEnabled != logsFileEnabled {
		tst.Errorf("Invalid logs.file-enabled - expected: %t, value: %t", logsFileEnabled, logs.FileEnabled)
	}

	if logs.MaxSize != logsMaxSize {
		tst.Errorf("Invalid logs.max-size - expected: %d, value: %d", logsMaxSize, logs.MaxSize)
	}

	if logs.MaxAge != logsMaxAge {
		tst.Errorf("Invalid logs.max-age - expected: %d, value: %d", logsMaxAge, logs.MaxAge)
	}
}

func TestLoadOtherStructSameProcess(tst *testing.T) {
	// Given
	msg := "application settings already loaded as different type"

	// When
	_, err := settings.Load[settings.Settings]()

	// Then
	if err == nil {
		tst.Errorf("Loading should end up with error")
	} else {
		if err.Error() != msg {
			tst.Errorf("Loading should end up with different error message - expected: %s, value: %s", msg, err.Error())
		}
	}
}
