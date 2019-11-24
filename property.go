package homie

import (
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

// Property represents a node's property according to the spec
type Property struct {
	Name     string   `mapstructure:"$name"`
	DataType DataType `mapstructure:"$datatype"`
	Value    string   `mapstructure:"_value"`
	Format   string   `mapstructure:"$format"`
	Unit     string   `mapstructure:"$unit"`
	Retained bool     `mapstructure:"$retained"`
	Settable bool     `mapstructure:"$settable"`
}

// NewProperty creates a new property
func NewProperty() *Property {
	return &Property{
		// Name:     id,
		Retained: true,
	}
}

// Publish publishes the property to MQTT at the given topic.
// It's the callers responsibility to ensure correct topic and include the
// property in the parent nodes $properties attribute.
func (p *Property) Publish(pub Publisher, topic string) {
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

func (p *Property) Unmarshal(subscribe Subscriber, topic string) {
	prefix := topic + "/"

	subscribe(prefix+"+", func(topic string, retained bool, message string) {
		topic = strings.TrimPrefix(topic, prefix)
		// fmt.Printf("prop: %s %v (%v)\n", topic, message, retained)

		switch topic {
		case "$name":
			p.Name = message
		case "$datatype":
			p.DataType = DataType(message)
		case "$format":
			p.Format = message
		case "$unit":
			p.Unit = message
		case "$retained":
			p.Retained = message == "true"
		case "$settable":
			p.Settable = message == "true"
		}
	})
}
