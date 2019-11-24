package homie

import (
	"testing"
)

func TestHomie(t *testing.T) {
	if h := NewHomie(); h.RootTopic != "homie" {
		t.Fail()
	}
}
func TestHomieAdd(t *testing.T) {
	h := NewHomie()
	d := NewDevice()

	if err := h.Add("dev", d); err != nil || h.Devices["dev"] != d {
		t.Fail()
	}

	if err := h.Add("dev", d); err == nil {
		t.Fail()
	}

	if d, err := h.NewDevice("dev2"); err != nil || h.Devices["dev2"] != d {
		t.Fail()
	}
}
