package homie

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"strings"
)

type DataType string

const (
	DataTypeString  DataType = "string"
	DataTypeInteger DataType = "integer"
	DataTypeFloat   DataType = "float"
	DataTypeBoolean DataType = "boolean"
	DataTypeEnum    DataType = "enum"
	DataTypeColor   DataType = "color"
)

type Property struct {
	ID       string
	Name     string   `mapstructure:"$name"`
	DataType DataType `mapstructure:"$datatype"`
	Value    string
	Format   string `mapstructure:"$format"`
	Unit     string `mapstructure:"$unit"`
	Retained bool   `mapstructure:"$retained"`
	Settable bool   `mapstructure:"$settable"`
}

func NewProperty(id string) *Property {
	return &Property{
		ID:       id,
		Name:     id,
		Retained: true,
	}
}

func (p *Property) Topic(base string) string {
	return base + "/" + p.ID
}

func (p *Property) Publish(pub Publisher, base string) {
	topic := p.Topic(base)

	// required attributes
	pub(topic+"/$name", true, p.Name)
	pub(topic+"/$datatype", true, string(p.DataType))

	// optional attributes
	if p.Format != "" {
		pub(topic+"/$format", true, p.Format)
	}
	if !p.Retained {
		pub(topic+"/$retained", true, "false")
	}
	if p.Settable {
		pub(topic+"/$settable", true, "true")
	}
	if p.Unit != "" {
		pub(topic+"/$unit", true, p.Unit)
	}

	pub(topic, p.Retained, p.Value)
}

// func (p *Property) Update(pub Publisher, base string) {
// 	topic := base + "/" + p.Name
// 	pub(topic, p.Retained, p.Value)
// }

func (p *Property) Unmarshal(subscribe Subscriber, base string) {
	prefix := p.Topic(base) + "/"

	subscribe(prefix+"+", func(topic string, retained bool, message string) {
		topic = strings.TrimPrefix(topic, prefix)
		fmt.Printf("prop: %s %v (%v)\n", topic, message, retained)

		// use mapstructure instead of decoding by property
		mapstructure.WeakDecode(map[string]string{
			topic: message,
		}, p)
	})
}
