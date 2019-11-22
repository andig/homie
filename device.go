package homie

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type State string

const (
	Version = "4.0.0"

	StateInit         State = "init"
	StateReady        State = "ready"
	StateDisconnected State = "disconnected"
	StateSleeping     State = "sleeping"
	StateLost         State = "lost"
	StateAlert        State = "alert"
)

type Device struct {
	ID             string
	Name           string `mapstructure:"$name"`
	State          State  `mapstructure:"$state"`
	Version        string `mapstructure:"$homie"`
	Implementation string `mapstructure:"$implementation"`
	Extensions     []string
	Nodes          map[string]*Node
}

func NewDevice(id string) *Device {
	return &Device{
		ID:         id,
		Name:       id,
		State:      StateInit,
		Version:    Version,
		Extensions: make([]string, 0),
		Nodes:      make(map[string]*Node, 0),
	}
}

func (d *Device) NewNode(id string) (*Node, error) {
	n := NewNode(id)
	return n, d.Add(n)
}

func (d *Device) Add(n *Node) error {
	if _, ok := d.Nodes[n.ID]; ok {
		return fmt.Errorf("node %s already exists", n.ID)
	}

	d.Nodes[n.ID] = n
	return nil
}

func (d *Device) Topic(base string) string {
	return base + "/" + d.ID
}

func (d *Device) Publish(pub Publisher, base string) {
	topic := d.Topic(base)

	pub(topic+"/$homie", true, d.Version)
	pub(topic+"/$name", true, d.Name)
	pub(topic+"/$state", true, string(d.State))
	pub(topic+"/$extensions", true, strings.Join(d.Extensions, ","))

	// optional attributes
	if d.Implementation != "" {
		pub(topic+"/$implementation", true, d.Implementation)
	}

	nodes := make([]string, 0, len(d.Nodes))
	for _, n := range d.Nodes {
		nodes = append(nodes, n.ID)
		n.Publish(pub, topic)
	}
	sort.Strings(nodes)
	pub(topic+"/$nodes", true, strings.Join(nodes, ","))
}

func (d *Device) Unmarshal(subscribe Subscriber, base string) {
	prefix := d.Topic(base) + "/"

	subscribe(prefix+"+", func(topic string, retained bool, message string) {
		topic = strings.TrimPrefix(topic, prefix)
		fmt.Printf("dev: %s %v (%v)\n", topic, message, retained)

		switch topic {
		case "$extensions":
			d.Extensions = strings.Split(message, ",")
		case "$nodes":
			nodes := strings.Split(message, ",")
			for _, id := range nodes {
				if _, ok := d.Nodes[id]; !ok {
					n, _ := d.NewNode(id)
					n.Unmarshal(subscribe, d.Topic(base))
				}
			}
		default:
			// use mapstructure instead of decoding by property
			mapstructure.WeakDecode(map[string]string{
				topic: message,
			}, d)
		}
	})
}
