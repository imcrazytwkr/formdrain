package form_config

type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeNumber  FieldType = "number"
	FieldTypeBoolean FieldType = "boolean"
	FieldTypeArray   FieldType = "array"
)

type FieldSchema struct {
	Fields []Field `json:"fields"`
}

type Field struct {
	Name     string      `json:"name"`
	Type     FieldType   `json:"type"`
	Required bool        `json:"required"`
	Items    *FieldItems `json:"items,omitempty"`
}

// FieldItems describes homogeneous array element types (string or number).
type FieldItems struct {
	Type FieldType `json:"type"`
}
