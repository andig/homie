package homie

import (
	"fmt"
	"sort"
	"strings"
)

const Version = "4.0.0"

type State string

const (
	StateInit         State = "init"
	StateReady        State = "ready"
	StateDisconnected State = "disconnected"
	StateSleeping     State = "sleeping"
	StateLost         State = "lost"
	StateAlert        State = "alert"
)

type Device struct {
	ID         string
	Name       string
	State      State
	Version    string
	Extensions []string
	Nodes      map[string]*Node
	Attributes map[string]string
}

func NewDevice(id string) *Device {
	return &Device{
		ID:         id,
		Name:       id,
		State:      StateInit,
		Version:    Version,
		Extensions: make([]string, 0),
		Nodes:      make(map[string]*Node, 0),
		Attributes: make(map[string]string, 0),
	}
}

func (d *Device) Add(n *Node) error {
	if _, ok := d.Nodes[n.ID]; ok {
		return fmt.Errorf("node %s already exists", n.ID)
	}

	d.Nodes[n.ID] = n
	return nil
}

func (d *Device) Publish(pub Publisher, base string) {
	topic := base + "/" + d.ID

	pub(topic+"/$name", true, d.Name)
	pub(topic+"/$version", true, d.Version)
	pub(topic+"/$state", true, string(d.State))
	pub(topic+"/$extensions", true, strings.Join(d.Extensions, ","))

	nodes := make([]string, 0, len(d.Nodes))
	for _, n := range d.Nodes {
		nodes = append(nodes, n.ID)
		n.Publish(pub, topic)
	}
	sort.Strings(nodes)
	pub(topic+"/$nodes", true, strings.Join(nodes, ","))
}
