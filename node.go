package homie

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type Node struct {
	ID         string
	Name       string `mapstructure:"$name"`
	Type       string `mapstructure:"$type"`
	Properties map[string]*Property
}

func NewNode(id string) *Node {
	return &Node{
		ID:         id,
		Name:       id,
		Properties: make(map[string]*Property, 0),
	}
}

func (n *Node) NewProperty(id string) (*Property, error) {
	p := NewProperty(id)
	return p, n.Add(p)
}

func (n *Node) Add(p *Property) error {
	if _, ok := n.Properties[p.ID]; ok {
		return fmt.Errorf("property %s already exists", p.ID)
	}

	n.Properties[p.ID] = p
	return nil
}

func (n *Node) Topic(base string) string {
	return base + "/" + n.ID
}

func (n *Node) Publish(pub Publisher, base string) {
	topic := n.Topic(base)

	pub(topic+"/$name", true, n.Name)
	pub(topic+"/$type", true, n.Type)

	properties := make([]string, 0, len(n.Properties))
	for _, prop := range n.Properties {
		properties = append(properties, prop.ID)
		prop.Publish(pub, topic)
	}
	sort.Strings(properties)
	pub(topic+"/$properties", true, strings.Join(properties, ","))
}

func (n *Node) Unmarshal(subscribe Subscriber, base string) {
	prefix := n.Topic(base) + "/"

	subscribe(prefix+"+", func(topic string, retained bool, message string) {
		topic = strings.TrimPrefix(topic, prefix)
		fmt.Printf("node: %s %v (%v)\n", topic, message, retained)

		switch topic {
		case "$properties":
			properties := strings.Split(message, ",")
			for _, id := range properties {
				if _, ok := n.Properties[id]; !ok {
					p, _ := n.NewProperty(id)
					p.Unmarshal(subscribe, n.Topic(base))
				}
			}
		default:
			// use mapstructure instead of decoding by property
			mapstructure.WeakDecode(map[string]string{
				topic: message,
			}, n)
		}
	})
}
