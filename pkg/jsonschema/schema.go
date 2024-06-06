package jsonschema

type Schema struct {
	Id                   string            `json:"id,omitempty"`
	Type                 SchemaType        `json:"type,omitempty"`
	Enum                 []string          `json:"enum,omitempty"`
	Minimum              *float64          `json:"minimum,omitempty"`
	ContentEncoding      string            `json:"content_encoding,omitempty"`
	Ref                  string            `json:"ref,omitempty"`
	Items                []Schema          `json:"items,omitempty"`
	AdditionalProperties *Schema           `json:"additional_properties,omitempty"`
	Properties           map[string]Schema `json:"properties,omitempty"`
}

type SchemaType string

const (
	ObjectType  SchemaType = "object"
	StringType  SchemaType = "string"
	NumberType  SchemaType = "number"
	IntegerType SchemaType = "integer"
	BooleanType SchemaType = "boolean"
	ArrayType   SchemaType = "array"
)
