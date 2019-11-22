package paho

import (
	"github.com/andig/homie"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Handler struct {
	Client mqtt.Client
	qos    byte
}

func NewHandler(deviceRoot string, template *mqtt.ClientOptions, qos byte) *Handler {
	opt := CloneOptions(template)
	opt.SetWill(deviceRoot+"/$state", string(homie.StateLost), qos, true)

	return &Handler{
		Client: mqtt.NewClient(opt),
		qos:    qos,
	}
}

func (h *Handler) Publish(topic string, retained bool, message string) {
	// fmt.Printf("%s %v (%v)\n", topic, message, retained)
	h.Client.Publish(topic, h.qos, retained, message)
}

func (h *Handler) Subscribe(topic string, callback homie.Callback) {
	// fmt.Printf("%s %v (%v)\n", topic, message, retained)
	h.Client.Subscribe(topic, h.qos, func(c mqtt.Client, m mqtt.Message) {
		callback(m.Topic(), m.Retained(), string(m.Payload()))
	})
}
