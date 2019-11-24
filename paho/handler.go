package paho

import (
	"errors"
	"time"

	"github.com/andig/homie"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// ErrorHandler is the type MQTT error handlers must implement
type ErrorHandler func(err error)

// Handler is the MQTT adapter for Eclipse Paho
type Handler struct {
	Client       mqtt.Client
	Timeout      time.Duration
	ErrorHandler ErrorHandler
	qos          byte
}

// NewHandler creates a new MQTT client connection with applied last will
func NewHandler(deviceRoot string, template *mqtt.ClientOptions, qos byte) *Handler {
	opt := CloneOptions(template)
	opt.SetWill(deviceRoot+"/$state", string(homie.StateLost), qos, true)

	return &Handler{
		Client: mqtt.NewClient(opt),
		qos:    qos,
	}
}

// handleToken invokes ErrorHandler if token times out or errors
func (h *Handler) handleToken(op string, token mqtt.Token) {
	if token.WaitTimeout(h.Timeout) {
		if e := token.Error(); e != nil {
			h.ErrorHandler(e)
		}
	} else {
		h.ErrorHandler(errors.New(op + " timeout"))
	}
}

// Publish implements homie.Publisher for Eclipse Paho
func (h *Handler) Publish(topic string, retained bool, message string) {
	// fmt.Printf("%s %v (%v)\n", topic, message, retained)
	token := h.Client.Publish(topic, h.qos, retained, message)
	if h.Timeout > 0 && h.ErrorHandler != nil {
		go h.handleToken("publish", token)
	}
}

// Subscribe implements homie.Subscriber for Eclipse Paho
func (h *Handler) Subscribe(topic string, callback homie.Callback) {
	// fmt.Printf("%s %v (%v)\n", topic, message, retained)
	token := h.Client.Subscribe(topic, h.qos, func(c mqtt.Client, m mqtt.Message) {
		callback(m.Topic(), m.Retained(), string(m.Payload()))
	})
	if h.Timeout > 0 && h.ErrorHandler != nil {
		go h.handleToken("subscribe", token)
	}
}
