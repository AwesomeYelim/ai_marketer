package llm

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
)

// GenerateSchema generates a JSON Schema from a Go type using reflection.
// It returns the schema as a raw JSON map suitable for use with the Claude API.
func GenerateSchema[T any]() map[string]interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)

	b, err := json.Marshal(schema)
	if err != nil {
		panic("failed to marshal JSON schema: " + err.Error())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		panic("failed to unmarshal JSON schema: " + err.Error())
	}

	return result
}

// SchemaToString converts a schema map to a formatted JSON string.
func SchemaToString[T any]() string {
	schema := GenerateSchema[T]()
	b, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}
