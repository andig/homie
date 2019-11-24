package main

import (
	"log"
	"time"

	"github.com/andig/homie"
	"github.com/andig/homie/paho"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	broker = "localhost:1883"
	qos    = byte(1)
)

func main() {
	// example device tree
	d := homie.NewDevice("meter")
	if n, _ := d.NewNode("tariff1"); true {
		n.Name = "Tarrif 1"

		if p, _ := n.NewProperty("energy"); true {
			p.Unit = "Wh"
			p.DataType = homie.DataTypeFloat
			p.Settable = true
		}
		if p, _ := n.NewProperty("power"); true {
			p.Unit = "W"
			p.DataType = homie.DataTypeFloat
		}
	}

	// template mqtt client options
	opt := mqtt.NewClientOptions()
	opt.AddBroker(broker)
	opt.SetAutoReconnect(true)

	// root topic for device
	topic := homie.DefaultRootTopic + "/meter"

	// mqtt client connection with cloned options and last will
	handler := paho.NewHandler(topic, opt, qos)
	handler.Timeout = 1 * time.Second
	handler.ErrorHandler = paho.Log
	if t := handler.Client.Connect(); !t.WaitTimeout(time.Second) {
		log.Fatalf("could not connect: %v", t.Error())
	}

	// publish the device using handler's Publish method
	d.Publish(handler.Publish, topic)
	time.Sleep(time.Second)

	// omitting the Disconnect() will set the device state to "lost"
	handler.Client.Disconnect(1000)
}
