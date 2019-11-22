package paho

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// CloneOptions creates a copy of the relevant mqtt options
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
