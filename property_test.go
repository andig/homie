package homie

import (
	"testing"
)

func TestProperty(t *testing.T) {
	if p := NewProperty("id"); p.ID != "id" || p.Name != "id" ||
		p.DataType != "" || p.Value != "" ||
		p.Format != "" || p.Unit != "" ||
		p.Settable != false || p.Retained != true {
		t.Fail()
	}
}

func TestPropertyPublish(t *testing.T) {
	p := NewProperty("id")
	p.Name = "name"
	p.DataType = DataTypeFloat
	p.Value = "1.23"
	p.Format = "1..2"
	p.Unit = "%"

	exp := []struct {
		t, m string
		r    bool
	}{
		{">/id/$name", "name", true},
		{">/id/$datatype", "float", true},
		{">/id/$format", "1..2", true},
		{">/id/$unit", "%", true},
		{">/id", "1.23", true},
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
	}, ">")

	if idx != len(exp) {
		t.Errorf("unexpected number of matches %d", idx)
	}
}

func TestPropertyPublishForUnretainedSettable(t *testing.T) {
	p := NewProperty("id")
	p.DataType = "float"
	p.Value = "1.23"
	p.Retained = false
	p.Settable = true

	exp := []struct {
		t, m string
		r    bool
	}{
		{">/id/$name", "id", true},
		{">/id/$datatype", "float", true},
		{">/id/$retained", "false", true},
		{">/id/$settable", "true", true},
		{">/id", "1.23", false},
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
	}, ">")

	if idx != len(exp) {
		t.Errorf("unexpected number of matches %d", idx)
	}
}
