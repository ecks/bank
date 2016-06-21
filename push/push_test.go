package push

import (
	"testing"

	"github.com/ksred/bank/configuration"
)

func TestSetConfig(t *testing.T) {
	// Load app config
	Config, err := configuration.LoadConfig()
	if err != nil {
		t.Errorf("TestPush.SetConfig: %v", err)
	}
	// Set config in packages
	SetConfig(&Config)
}
