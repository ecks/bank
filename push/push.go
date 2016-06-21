// Package push sends push notifications to devices
// These push notifications are linked to users accounts
// through the accounts_push_tokens table
package push

import (
	"errors"

	"github.com/ksred/apns"
	"github.com/ksred/bank/configuration"
)

var Config configuration.Configuration

type PushDevice struct {
	Token    string
	Platform string
}

func SetConfig(config *configuration.Configuration) {
	Config = *config
}

func SendNotification(accountNumber string, message string, badge uint, sound string) (err error) {
	// Get any push tokens for the user
	pushDevices, err := getPushTokens(accountNumber)
	if err != nil {
		return errors.New("push.SendNotification: Could not get push devices " + err.Error())
	}
	// Loop through
	for _, pd := range pushDevices {
		// Switch on device type
		switch pd.Platform {
		case "ios":
			doSendNotificationAPNS(pd.Token, message, badge, sound)
		default:
			// Not supported yet
			return nil
		}
	}
	return nil
}

func doSendNotificationAPNS(token string, message string, badge uint, sound string) {
	// Set vars
	var gateway string
	var apnsCert string
	var apnsKey string
	var notificationSound string

	// Get env
	switch Config.PushEnv {
	case "development", "production":
		gateway = apns.ProductionGateway
		apnsCert = "../certs/apns-cert.pem"
		apnsKey = "../certs/apns-prod.pem"
		break
	}

	// Fetch relevant sound
	switch sound {
	case "alert", "default":
		sound = "bvnk_default.aiff"
		break
	}

	// Send notification
	c, _ := apns.NewClient(gateway, apnsCert, apnsKey)

	p := apns.NewPayload()
	p.APS.Alert.Body = message
	p.APS.Badge.Set(badge)
	p.APS.Sound = notificationSound

	m := apns.NewNotification()
	m.Payload = p
	m.DeviceToken = token
	m.Priority = apns.PriorityImmediate

	c.Send(m)
}
