package homie

import (
	"strings"
	"testing"
)

func TestDevice(t *testing.T) {
	if d := NewDevice(); d.State != StateInit || d.Version != Version || d.Implementation != "" {
		t.Fail()
	}
}
func TestDeviceAdd(t *testing.T) {
	d := NewDevice()
	n := NewNode()

	if err := d.Add("node", n); err != nil || d.Nodes["node"] != n {
		t.Fail()
	}

	if err := d.Add("node", n); err == nil {
		t.Fail()
	}

	if n, err := d.NewNode("node2"); err != nil || d.Nodes["node2"] != n {
		t.Fail()
	}
}

func TestDevicePublish(t *testing.T) {
	d := NewDevice()
	d.Name = "name"
	d.Implementation = "impl"
	d.State = StateReady
	d.Extensions = append(d.Extensions, "foo", "bar")
	d.NewNode("n1")
	d.NewNode("n2")

	exp := []struct {
		t, m string
		r    bool
	}{
		{"homie/dev/$homie", Version, true},
		{"homie/dev/$name", "name", true},
		{"homie/dev/$state", string(StateReady), true},
		{"homie/dev/$extensions", "foo,bar", true},
		{"homie/dev/$implementation", "impl", true},
		{"homie/dev/$nodes", "n1,n2", true},
	}

	idx := 0
	d.Publish(func(topic string, retained bool, message string) {
		// fmt.Printf("%s %v (%v)\n", topic, message, retained)
		// filter node properties
		if strings.Contains(topic, "/n1/") || strings.Contains(topic, "/n2/") {
			return
		}

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
	}, "homie/dev")

	if idx != len(exp) {
		t.Errorf("unexpected number of matches %d", idx)
	}
}
