package homie

import (
	"strings"
	"testing"
)

func TestNode(t *testing.T) {
	if n := NewNode(); n.Type != "" {
		t.Fail()
	}
}
func TestNodeAdd(t *testing.T) {
	n := NewNode()
	p := NewProperty()

	if err := n.Add("prop", p); err != nil || n.Properties["prop"] != p {
		t.Fail()
	}

	if err := n.Add("prop", p); err == nil {
		t.Fail()
	}

	if p, err := n.NewProperty("prop2"); err != nil || n.Properties["prop2"] != p {
		t.Fail()
	}
}

func TestNodePublish(t *testing.T) {
	n := NewNode()
	n.Name = "name"
	n.Type = "type"
	n.NewProperty("p1")
	n.NewProperty("p2")

	exp := []struct {
		t, m string
		r    bool
	}{
		{"homie/dev/node/$name", "name", true},
		{"homie/dev/node/$type", "type", true},
		{"homie/dev/node/$properties", "p1,p2", true},
	}

	idx := 0
	n.Publish(func(topic string, retained bool, message string) {
		// fmt.Printf("%s %v (%v)\n", topic, message, retained)
		// filter node properties
		if strings.Contains(topic, "/p1") || strings.Contains(topic, "/p2") {
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
	}, "homie/dev/node")

	if idx != len(exp) {
		t.Errorf("unexpected number of matches %d", idx)
	}
}
