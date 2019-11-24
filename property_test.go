package homie

import (
	"testing"
)

func TestProperty(t *testing.T) {
	if p := NewProperty(); p.DataType != "" || p.Value != "" ||
		p.Format != "" || p.Unit != "" ||
		p.Settable != false || p.Retained != true {
		t.Fail()
	}
}

func TestPropertyPublish(t *testing.T) {
	p := NewProperty()
	p.Name = "name"
	p.DataType = DataTypeFloat
	p.Value = "1.23"
	p.Format = "1..2"
	p.Unit = "%"

	exp := []struct {
		t, m string
		r    bool
	}{
		{"homie/dev/node/prop/$name", "name", true},
		{"homie/dev/node/prop/$datatype", "float", true},
		{"homie/dev/node/prop/$format", "1..2", true},
		{"homie/dev/node/prop/$unit", "%", true},
		{"homie/dev/node/prop", "1.23", true},
	}

	idx := 0
	p.Publish(func(topic string, retained bool, message string) {
		// fmt.Printf("%s %v (%v)\n", topic, message, retained)
		if idx >= len(exp) {
			t.Errorf("unexpected index %d", idx)
			return
		}

		e := exp[idx]
		if e.t != topic || e.m != message || e.r != retained {
			t.Errorf("expected %s %s %v", e.t, e.m, e.r)
			t.Errorf("got %s %s %v", topic, message, retained)
		}
		idx++
	}, "homie/dev/node/prop")

	if idx != len(exp) {
		t.Errorf("unexpected number of matches %d", idx)
	}
}

func TestPropertyPublishForUnretainedSettable(t *testing.T) {
	p := NewProperty()
	p.Name = "name"
	p.DataType = "float"
	p.Value = "1.23"
	p.Retained = false
	p.Settable = true

	exp := []struct {
		t, m string
		r    bool
	}{
		{"homie/dev/node/prop/$name", "name", true},
		{"homie/dev/node/prop/$datatype", "float", true},
		{"homie/dev/node/prop/$retained", "false", true},
		{"homie/dev/node/prop/$settable", "true", true},
		{"homie/dev/node/prop", "1.23", false},
	}

	idx := 0
	p.Publish(func(topic string, retained bool, message string) {
		// fmt.Printf("%s %v (%v)\n", topic, message, retained)
		if idx >= len(exp) {
			t.Errorf("unexpected index %d", idx)
			return
		}

		e := exp[idx]
		if e.t != topic || e.m != message || e.r != retained {
			t.Errorf("expected %s %s %v", e.t, e.m, e.r)
			t.Errorf("got %s %s %v", topic, message, retained)
		}
		idx++
	}, "homie/dev/node/prop")

	if idx != len(exp) {
		t.Errorf("unexpected number of matches %d", idx)
	}
}
