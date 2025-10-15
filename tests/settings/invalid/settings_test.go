package invalid

import (
	"testing"

	"github.com/ternaryss/rest2go/pkg/rest2go/settings"
)

type AppSettings struct {
	Foo string `yaml:"foo"`
}

func TestLoadWithoutDefaults(tst *testing.T) {
	// Given
	msg := "defaults not implemented for application settings structure"

	// When
	_, err := settings.Load[AppSettings]()

	// Then
	if err == nil {
		tst.Errorf("Loading should end up with error")
	} else {
		if err.Error() != msg {
			tst.Errorf("Loading should end up with different error message - expected: %s, value: %s", msg, err.Error())
		}
	}
}
