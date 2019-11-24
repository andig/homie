package homie

import (
	"fmt"
	"sort"
	"strings"
)

// State is the Device state
type State string

const (
	// Version is the spec version
	Version = "4.0.0"

	StateInit         State = "init"
	StateReady        State = "ready"
	StateDisconnected State = "disconnected"
	StateSleeping     State = "sleeping"
	StateLost         State = "lost"
	StateAlert        State = "alert"
)

// Device represents a device in terms of the spec
type Device struct {
	Name           string           `mapstructure:"$name"`
	State          State            `mapstructure:"$state"`
	Version        string           `mapstructure:"$homie"`
	Implementation string           `mapstructure:"$implementation"`
	Extensions     []string         `mapstructure:"_extensions"`
	Nodes          map[string]*Node `mapstructure:"_nodes"`
}

// NewDevice creates a new device
func NewDevice() *Device {
	return &Device{
		// Name:       id,
		State:      StateInit,
		Version:    Version,
		Extensions: make([]string, 0),
		Nodes:      make(map[string]*Node),
	}
}

// NewNode is a conveniece method for creating a new node and attaching it to the device
func (d *Device) NewNode(id string) (*Node, error) {
	n := NewNode()
	return n, d.Add(id, n)
}

// Add attaches a node to the node. An error is raised on duplicate node id.
func (d *Device) Add(id string, n *Node) error {
	if _, ok := d.Nodes[id]; ok {
		return fmt.Errorf("node %s already exists", id)
	}

	d.Nodes[id] = n
	return nil
}

// Publish publishes the device including nodes/properties to MQTT at the given topic.
// It's the callers responsibility to ensure correct topic.
func (d *Device) Publish(pub Publisher, topic string) {
	pub(topic+"/$homie", true, d.Version)
	pub(topic+"/$name", true, d.Name)
	pub(topic+"/$state", true, string(d.State))
	pub(topic+"/$extensions", true, strings.Join(d.Extensions, ","))

	// optional attributes
	if d.Implementation != "" {
		pub(topic+"/$implementation", true, d.Implementation)
	}

	nodes := make([]string, 0, len(d.Nodes))
	for id, n := range d.Nodes {
		nodes = append(nodes, id)
		n.Publish(pub, topic+"/"+id)
	}
	sort.Strings(nodes)
	pub(topic+"/$nodes", true, strings.Join(nodes, ","))
}

func (d *Device) Unmarshal(subscribe Subscriber, topic string) {
	prefix := topic + "/"

	subscribe(prefix+"+", func(topic string, retained bool, message string) {
		topic = strings.TrimPrefix(topic, prefix)
		// fmt.Printf("dev: %s %v (%v)\n", topic, message, retained)

		switch topic {
		case "$homie":
			d.Version = message
		case "$name":
			d.Name = message
		case "$state":
			d.State = State(message)
		case "$implementation":
			d.Implementation = message
		case "$extensions":
			d.Extensions = strings.Split(message, ",")
		case "$nodes":
			nodes := strings.Split(message, ",")
			for _, id := range nodes {
				if _, ok := d.Nodes[id]; !ok {
					n, _ := d.NewNode(id)
					n.Unmarshal(subscribe, prefix+id)
				}
			}
		}
	})
}
