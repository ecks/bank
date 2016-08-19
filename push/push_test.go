package push

import (
	"testing"

	"github.com/bvnk/bank/configuration"
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

func TestDoSendNotificationAPNS(t *testing.T) {
	err := doSendNotificationAPNS("TOKEN", "Testing from CLI ⚡️", 1, "default")
	if err != nil {
		t.Errorf("TestPush.SendNotification: %v", err)
	}
}
