package jsonschema

const DefinitionsPath = "#/definitions/"
const SchemaVersion = "http://json-schema.org/draft-07/schema#"

type Schema struct {
	SchemaVersion        string             `json:"$schema,omitempty"`
	Id                   string             `json:"$id,omitempty"`
	Type                 SchemaType         `json:"type,omitempty"`
	Enum                 []string           `json:"enum,omitempty"`
	Minimum              *float64           `json:"minimum,omitempty"`
	ContentEncoding      string             `json:"contentEncoding,omitempty"`
	Ref                  string             `json:"$ref,omitempty"`
	Items                []*Schema          `json:"items,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	Definitions          map[string]*Schema `json:"definitions,omitempty"`
	AdditionalProperties bool               `json:"additionalProperties,omitempty"`
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

func (receiver SchemaType) String() string {
	return string(receiver)
}
