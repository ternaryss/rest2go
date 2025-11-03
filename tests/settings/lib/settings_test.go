package lib

import (
	"testing"

	"github.com/ternaryss/rest2go/pkg/rest2go/settings"
)

func TestLoadDefaults(tst *testing.T) {
	// Given
	logsLevel := "info"
	logsFileEnabled := false
	logsMaxSize := 10
	logsMaxAge := 7

	// When
	settings, err := settings.Load[settings.Settings]()

	// Then
	if err != nil {
		tst.Errorf("Loading default settings failed: %s", err)
	}

	logs := settings.Logs

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
