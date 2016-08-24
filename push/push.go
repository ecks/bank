// Package push sends push notifications to devices
// These push notifications are linked to users accounts
// through the accounts_push_tokens table
package push

import (
	"errors"
	"fmt"

	"github.com/ksred/apns"
)

type PushDevice struct {
	Token    string
	Platform string
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
			err = doSendNotificationAPNS(pd.Token, message, badge, sound)
			if err != nil {
				return err
			}
		default:
			// Not supported yet
			return nil
		}
	}
	return nil
}

func doSendNotificationAPNS(token string, message string, badge uint, sound string) (err error) {
	// Set vars
	var gateway string
	var apnsCert string
	var apnsKey string
	//var notificationSound string

	// Get env
	switch Config.PushEnv {
	case "development", "production":
		gateway = apns.ProductionGateway
		apnsCert = Config.ApplePushCert
		apnsKey = Config.ApplePushKey
		break
	}

	// Fetch relevant sound
	/*
		switch sound {
		case "alert", "default":
			sound = "bvnk_default.aiff"
			break
		}
	*/

	// Send notification
	c, err := apns.NewClientWithFiles(gateway, apnsCert, apnsKey)
	if err != nil {
		return errors.New("Could not create client: " + err.Error())
	}
	fmt.Println(gateway)
	fmt.Println(apnsCert)
	fmt.Println(apnsKey)

	p := apns.NewPayload()
	p.APS.Alert.Body = message
	p.APS.Badge.Set(badge)

	//p.SetCustomValue("link", "yourapp://precache/20140718")

	m := apns.NewNotification()
	m.Payload = p
	m.DeviceToken = token
	m.Priority = apns.PriorityImmediate
	m.Identifier = uint32(1)

	c.Send(m)

	return
}
