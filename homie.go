package homie

import "fmt"

type Homie struct {
	base    string
	Devices map[string]*Device
}

type Publisher func(topic string, retained bool, message string)

const Base = "homie"

func NewHomie(base ...string) *Homie {
	topic := Base
	if len(base) == 1 {
		topic = base[0]
	}
	return &Homie{
		base:    topic,
		Devices: make(map[string]*Device, 0),
	}
}

func (h *Homie) Add(d *Device) error {
	if _, ok := h.Devices[d.ID]; ok {
		return fmt.Errorf("device %s already exists", d.ID)
	}

	h.Devices[d.ID] = d
	return nil
}

func (h *Homie) Publish(p Publisher) {
	for _, d := range h.Devices {
		d.Publish(p, h.base)
	}
}
