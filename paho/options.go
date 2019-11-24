package paho

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// CloneOptions creates a copy of the relevant mqtt options.
// It is used to create a per-device MQTT connection from template options.
// This is needed to satisfy the Homie spec which requires a dedicated connection
// to apply the last will to.
func CloneOptions(o *mqtt.ClientOptions) *mqtt.ClientOptions {
	opt := mqtt.NewClientOptions()

	opt.SetUsername(o.Username)
	opt.SetPassword(o.Password)
	opt.SetClientID(o.ClientID)
	opt.SetCleanSession(o.CleanSession)
	opt.SetAutoReconnect(o.AutoReconnect)

	for _, b := range o.Servers {
		opt.AddBroker(b.String())
	}

	return opt
}
