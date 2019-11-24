package homie

import (
	"fmt"
	"sort"
	"strings"
)

// Node represents a device's node in terms of the spec
type Node struct {
	Name       string               `mapstructure:"$name"`
	Type       string               `mapstructure:"$type"`
	Properties map[string]*Property `mapstructure:"_properties"`
}

// NewNode creates a new node
func NewNode() *Node {
	return &Node{
		// Name:       id,
		Properties: make(map[string]*Property),
	}
}

// NewProperty is a conveniece method for creating a new property and attaching it to the nodes
func (n *Node) NewProperty(id string) (*Property, error) {
	p := NewProperty()
	return p, n.Add(id, p)
}

// Add attaches a property to the node. An error is raised on duplicate property id.
func (n *Node) Add(id string, p *Property) error {
	if _, ok := n.Properties[id]; ok {
		return fmt.Errorf("property %s already exists", id)
	}

	n.Properties[id] = p
	return nil
}

// Publish publishes the node including properties to MQTT at the given topic.
// It's the callers responsibility to ensure correct topic and include node
// in the parent devices $nodes attribute.
func (n *Node) Publish(pub Publisher, topic string) {
	pub(topic+"/$name", true, n.Name)
	pub(topic+"/$type", true, n.Type)

	properties := make([]string, 0, len(n.Properties))
	for id, p := range n.Properties {
		properties = append(properties, id)
		p.Publish(pub, topic+"/"+id)
	}
	sort.Strings(properties)
	pub(topic+"/$properties", true, strings.Join(properties, ","))
}

func (n *Node) Unmarshal(subscribe Subscriber, topic string) {
	prefix := topic + "/"

	subscribe(prefix+"+", func(topic string, retained bool, message string) {
		topic = strings.TrimPrefix(topic, prefix)
		// fmt.Printf("node: %s %v (%v)\n", topic, message, retained)

		switch topic {
		case "$name":
			n.Name = message
		case "$type":
			n.Type = message
		case "$properties":
			properties := strings.Split(message, ",")
			for _, id := range properties {
				if _, ok := n.Properties[id]; !ok {
					p, _ := n.NewProperty(id)
					p.Unmarshal(subscribe, prefix+id)
				}
			}
		}
	})
}
