package homie

import (
	"fmt"
	"strings"
)

type Homie struct {
	RootTopic string
	Devices   map[string]*Device
}

type Publisher func(topic string, retained bool, message string)
type Subscriber func(topic string, callback Callback)
type Callback func(topic string, retained bool, message string)

const DefaultRootTopic = "homie"

func NewHomie(rootTopic ...string) *Homie {
	topic := DefaultRootTopic
	if len(rootTopic) == 1 {
		topic = rootTopic[0]
	}
	return &Homie{
		RootTopic: topic,
		Devices:   make(map[string]*Device, 0),
	}
}

func (h *Homie) NewDevice(id string) (*Device, error) {
	d := NewDevice(id)
	return d, h.Add(d)
}

func (h *Homie) Add(d *Device) error {
	if _, ok := h.Devices[d.ID]; ok {
		return fmt.Errorf("device %s already exists", d.ID)
	}

	h.Devices[d.ID] = d
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
