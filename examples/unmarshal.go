package main

import (
	"fmt"
	"log"
	"time"

	"github.com/andig/homie"
	"github.com/andig/homie/paho"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	root   = homie.DefaultRootTopic // "homie"
	broker = "localhost:1883"
	qos    = byte(1)
)

func main() {
	// example device tree
	h := homie.NewHomie(root)

	// template mqtt client options
	opt := mqtt.NewClientOptions()
	opt.AddBroker(broker)

	// mqtt client connection with cloned options and last will
	handler := paho.NewHandler("", opt, qos)
	if t := handler.Client.Connect(); !t.WaitTimeout(time.Second) {
		log.Fatalf("could not connect: %v", t.Error())
	}

	// unmarshal the device tree
	h.Unmarshal(handler.Subscribe)
	time.Sleep(time.Second)

	// stop further changes to the device tree
	handler.Client.Disconnect(0)

	// print the unmarshaled device tree
	for _, d := range h.Devices {
		d.Publish(func(topic string, retained bool, message string) {
			fmt.Printf("%s \"%v\" (%v)\n", topic, message, retained)
		}, root)
	}
}
