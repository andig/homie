package homie

import (
	"strings"
	"testing"
)

func TestNode(t *testing.T) {
	if n := NewNode("id"); n.ID != "id" || n.Name != "id" || n.Type != "" {
		t.Fail()
	}
}
func TestNodeAdd(t *testing.T) {
	n := NewNode("id")
	p := NewProperty("prop")

	if err := n.Add(p); err != nil || n.Properties["prop"] != p {
		t.Fail()
	}

	if err := n.Add(p); err == nil {
		t.Fail()
	}

	if p, err := n.NewProperty("prop2"); err != nil || n.Properties["prop2"] != p {
		t.Fail()
	}
}

func TestNodePublish(t *testing.T) {
	n := NewNode("id")
	n.Type = "type"
	n.Add(NewProperty("p1"))
	n.Add(NewProperty("p2"))

	exp := []struct {
		t, m string
		r    bool
	}{
		{">/id/$name", "id", true},
		{">/id/$type", "type", true},
		{">/id/$properties", "p1,p2", true},
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
	}, ">")

	if idx != len(exp) {
		t.Errorf("unexpected number of matches %d", idx)
	}
}
