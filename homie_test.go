package homie

import (
	"testing"
)

func TestHomie(t *testing.T) {
	if h := NewHomie(); h.base != "homie" {
		t.Fail()
	}
	if h := NewHomie("alternate"); h.base != "alternate" {
		t.Fail()
	}
}
func TestHomieAdd(t *testing.T) {
	h := NewHomie()
	d := NewDevice("dev")

	if err := h.Add(d); err != nil || h.Devices["dev"] != d {
		t.Fail()
	}

	if err := h.Add(d); err == nil {
		t.Fail()
	}
}
