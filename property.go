package homie

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
	Name     string
	DataType DataType
	Value    string
	Format   string
	Unit     string
	Retained bool
	Settable bool
}

func NewProperty(id string) *Property {
	return &Property{
		ID:       id,
		Name:     id,
		Retained: true,
	}
}

func (p *Property) Publish(pub Publisher, base string) {
	topic := base + "/" + p.ID

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
