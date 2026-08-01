package form_config

type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeNumber  FieldType = "number"
	FieldTypeBoolean FieldType = "boolean"
	FieldTypeArray   FieldType = "array"
)

type FieldSchema struct {
	Version int     `json:"version"`
	Fields  []Field `json:"fields"`
}

type Field struct {
	Name     string      `json:"name"`
	Type     FieldType   `json:"type"`
	Required bool        `json:"required"`
	Items    *FieldItems `json:"items,omitempty"`
}

// FieldItems describes homogeneous array element types (MVP: string|number|boolean).
type FieldItems struct {
	Type FieldType `json:"type"`
}
