package homie

import (
	"fmt"
	"sort"
	"strings"
)

type Node struct {
	ID         string
	Name       string
	Type       string
	Properties map[string]*Property
}

func NewNode(id string) *Node {
	return &Node{
		ID:         id,
		Name:       id,
		Properties: make(map[string]*Property, 0),
	}
}

func (n *Node) Add(p *Property) error {
	if _, ok := n.Properties[p.ID]; ok {
		return fmt.Errorf("property %s already exists", p.ID)
	}

	n.Properties[p.ID] = p
	return nil
}

func (n *Node) Publish(pub Publisher, base string) {
	topic := base + "/" + n.ID

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
