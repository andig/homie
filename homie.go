package homie

import (
	"fmt"
	"strings"
)

type Homie struct {
	RootTopic string             `mapstructure:"_roottopic"`
	Devices   map[string]*Device `mapstructure:"_devices"`
}

type Publisher func(topic string, retained bool, message string)
type Subscriber func(topic string, callback Callback)
type Callback func(topic string, retained bool, message string)

const DefaultRootTopic = "homie"

func NewHomie() *Homie {
	return &Homie{
		RootTopic: DefaultRootTopic,
		Devices:   make(map[string]*Device),
	}
}

func (h *Homie) NewDevice(id string) (*Device, error) {
	d := NewDevice()
	return d, h.Add(id, d)
}

func (h *Homie) Add(id string, d *Device) error {
	if _, ok := h.Devices[id]; ok {
		return fmt.Errorf("device %s already exists", id)
	}

	h.Devices[id] = d
	return nil
}

func (h *Homie) Unmarshal(subscribe Subscriber) {
	prefix := h.RootTopic + "/"
	subscribe(prefix+"+/+", func(topic string, retained bool, message string) {
		topic = strings.TrimPrefix(topic, prefix)
		fmt.Printf("homie: %s %v (%v)\n", topic, message, retained)

		segments := strings.Split(topic, "/")
		id := segments[0]

		if _, ok := h.Devices[id]; !ok {
			d, _ := h.NewDevice(id)
			d.Unmarshal(subscribe, h.RootTopic)
		}
	})
}
