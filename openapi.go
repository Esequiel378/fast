package fast

import (
	"encoding/json"
	"path"
	"reflect"
	"slices"
	"strings"
)

// OpenAPIInfo contains basic information about the API
type OpenAPIInfo struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
}

// OpenAPISchema represents the OpenAPI schema specification
type OpenAPISchema struct {
	OpenAPI    string                    `json:"openapi"`
	Info       OpenAPIInfo               `json:"info"`
	Paths      map[string]PathItemObject `json:"paths"`
	Components ComponentsObject          `json:"components"`
	Tags       []TagObject               `json:"tags,omitempty"` // Added Tags field
}

// TagObject represents an OpenAPI tag
type TagObject struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// PathItemObject holds the operations for a specific path
type PathItemObject map[string]OperationObject

// OperationObject describes a single API operation on a path
type OperationObject struct {
	Summary     string                    `json:"summary,omitempty"`
	Description string                    `json:"description,omitempty"`
	OperationID string                    `json:"operationId,omitempty"`
	Tags        []string                  `json:"tags,omitempty"`
	Parameters  []ParameterObject         `json:"parameters,omitempty"`
	RequestBody *RequestBodyObject        `json:"requestBody,omitempty"`
	Responses   map[string]ResponseObject `json:"responses"`
}

// ParameterObject describes a single operation parameter
type ParameterObject struct {
	Name        string       `json:"name"`
	In          string       `json:"in"` // query, header, path, cookie
	Description string       `json:"description,omitempty"`
	Required    bool         `json:"required"`
	Schema      SchemaObject `json:"schema"`
}

// RequestBodyObject describes a request body
type RequestBodyObject struct {
	Description string                     `json:"description,omitempty"`
	Content     map[string]MediaTypeObject `json:"content"`
	Required    bool                       `json:"required,omitempty"`
}

// MediaTypeObject provides schema for the media type
type MediaTypeObject struct {
	Schema SchemaObject `json:"schema"`
}

// ResponseObject describes a single response from an API operation
type ResponseObject struct {
	Description string                     `json:"description"`
	Content     map[string]MediaTypeObject `json:"content,omitempty"`
}

// SchemaObject describes the object schema
type SchemaObject struct {
	Type       string                  `json:"type,omitempty"`
	Format     string                  `json:"format,omitempty"`
	Properties map[string]SchemaObject `json:"properties,omitempty"`
	Items      *SchemaObject           `json:"items,omitempty"`
	Ref        string                  `json:"$ref,omitempty"`
	Required   []string                `json:"required,omitempty"`
}

// ComponentsObject holds schemas that can be reused
type ComponentsObject struct {
	Schemas map[string]SchemaObject `json:"schemas,omitempty"`
}

// OpenAPIGenerator is responsible for creating OpenAPI documentation
type OpenAPIGenerator struct {
	handlers     map[string]Handler
	info         OpenAPIInfo
	schemas      map[string]SchemaObject
	tagsByName   map[string]TagObject // Map to store unique tags
	tagsForPaths map[string][]string  // Store tag associations for paths
}

// NewOpenAPIGenerator creates a new instance of OpenAPIGenerator
func NewOpenAPIGenerator(info OpenAPIInfo) *OpenAPIGenerator {
	return &OpenAPIGenerator{
		handlers:     make(map[string]Handler),
		info:         info,
		schemas:      make(map[string]SchemaObject),
		tagsByName:   make(map[string]TagObject),
		tagsForPaths: make(map[string][]string),
	}
}

// RegisterHandler adds a handler to be documented
func (g *OpenAPIGenerator) RegisterHandler(rootPath string, handler Handler) {
	path := path.Join(rootPath, handler.Path())
	g.handlers[path] = handler

	// Auto-generate tag for this path
	g.generateTagsForPath(path)
}

// generateTagsForPath extracts meaningful tags from a path
func (g *OpenAPIGenerator) generateTagsForPath(pathStr string) {
	// Handle empty path
	if pathStr == "" {
		return
	}

	// Split the path into segments
	segments := strings.Split(strings.Trim(pathStr, "/"), "/")

	if len(segments) == 0 {
		return
	}

	// Use the first segment as the primary tag
	primaryTag := segments[0]
	if primaryTag == "" {
		return
	}

	// Convert to title case for nicer display
	tagName := toTitleCase(primaryTag)

	// Add to tags map if not exists
	if _, exists := g.tagsByName[tagName]; !exists {
		g.tagsByName[tagName] = TagObject{
			Name:        tagName,
			Description: "Operations related to " + tagName,
		}
	}

	// Store the association between path and tags
	g.tagsForPaths[pathStr] = []string{tagName}

	// If path has a second segment, consider it a sub-resource
	if len(segments) > 1 && segments[1] != "" {
		// For paths like /admin/users, we might want a secondary tag "Admin Users"
		subTag := toTitleCase(primaryTag + " " + segments[1])

		if _, exists := g.tagsByName[subTag]; !exists {
			g.tagsByName[subTag] = TagObject{
				Name:        subTag,
				Description: "Operations related to " + subTag,
			}
		}

		// Add secondary tag
		g.tagsForPaths[pathStr] = append(g.tagsForPaths[pathStr], subTag)
	}
}

// toTitleCase converts a string to title case (e.g., "api-users" -> "Api Users")
func toTitleCase(s string) string {
	// Replace hyphens and underscores with spaces
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")

	// Split into words
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			// Capitalize first letter
			words[i] = strings.ToUpper(word[0:1]) + word[1:]
		}
	}

	return strings.Join(words, " ")
}

// GenerateSchema generates the OpenAPI schema for all registered handlers
func (g *OpenAPIGenerator) GenerateSchema() (*OpenAPISchema, error) {
	schema := &OpenAPISchema{
		OpenAPI: "3.0.3",
		Info:    g.info,
		Paths:   make(map[string]PathItemObject),
		Components: ComponentsObject{
			Schemas: make(map[string]SchemaObject),
		},
	}

	// Process each handler to build paths
	for path, handler := range g.handlers {
		g.processHandler(path, schema, handler)
	}

	// Add collected schemas to components
	schema.Components.Schemas = g.schemas

	// Convert tags map to slice for the OpenAPI schema
	for _, tag := range g.tagsByName {
		schema.Tags = append(schema.Tags, tag)
	}

	return schema, nil
}

// processHandler processes a single handler to extract path, method, and schemas
func (g *OpenAPIGenerator) processHandler(path string, schema *OpenAPISchema, handler Handler) {
	method := strings.ToLower(handler.Method())

	var (
		inputType  = reflect.TypeOf(handler.InputSerializer())
		outputType = reflect.TypeOf(handler.OutputSerializer())
	)

	// Create or update path
	if _, exists := schema.Paths[path]; !exists {
		schema.Paths[path] = make(PathItemObject)
	}

	// Create operation
	operation := OperationObject{
		OperationID: method + strings.ReplaceAll(path, "/", "_"),
		Responses:   make(map[string]ResponseObject),
	}

	// Add tags to the operation
	if tags, exists := g.tagsForPaths[path]; exists && len(tags) > 0 {
		operation.Tags = tags
	}

	// Add request body for non-GET methods
	if method != "get" && inputType != nil {
		inputSchema := g.generateSchemaForType(inputType)
		inputName := inputType.Name()
		if inputName != "" && inputName != "In" {
			g.schemas[inputName] = inputSchema
			operation.RequestBody = &RequestBodyObject{
				Content: map[string]MediaTypeObject{
					"application/json": {
						Schema: SchemaObject{
							Ref: "#/components/schemas/" + inputName,
						},
					},
				},
				Required: true,
			}
		}
	} else if inputType != nil && method == "get" {
		// For GET requests, add parameters from the input type
		operation.Parameters = g.generateParametersForType(inputType)
	}

	// Add response
	if outputType != nil {
		outputSchema := g.generateSchemaForType(outputType)
		outputName := outputType.Name()
		if outputName != "" && outputName != "Out" {
			g.schemas[outputName] = outputSchema
			operation.Responses["200"] = ResponseObject{
				Description: "Successful operation",
				Content: map[string]MediaTypeObject{
					"application/json": {
						Schema: SchemaObject{
							Ref: "#/components/schemas/" + outputName,
						},
					},
				},
			}
		} else {
			// Default output response
			operation.Responses["200"] = ResponseObject{
				Description: "Successful operation",
				Content: map[string]MediaTypeObject{
					"application/json": {
						Schema: outputSchema,
					},
				},
			}
		}
	} else {
		// Default response if no output type is found
		operation.Responses["200"] = ResponseObject{
			Description: "Successful operation",
		}
	}

	// Add error responses
	operation.Responses["400"] = ResponseObject{
		Description: "Bad request",
	}
	operation.Responses["422"] = ResponseObject{
		Description: "Validation error",
	}
	operation.Responses["500"] = ResponseObject{
		Description: "Internal server error",
	}

	schema.Paths[path][method] = operation
}

// generateSchemaForType generates an OpenAPI schema for a Go type
func (g *OpenAPIGenerator) generateSchemaForType(t reflect.Type) SchemaObject {
	if t == nil {
		return SchemaObject{Type: "object"}
	}

	// Dereference pointer if needed
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Struct:
		schema := SchemaObject{
			Type:       "object",
			Properties: make(map[string]SchemaObject),
		}

		var required []string

		for i := range t.NumField() {
			field := t.Field(i)

			// Skip unexported fields
			if !field.IsExported() {
				continue
			}

			// Get JSON tag
			tag := field.Tag.Get("json")
			if tag == "-" {
				continue
			}

			// Parse the tag
			parts := strings.Split(tag, ",")
			name := parts[0]
			if name == "" {
				name = field.Name
			}

			// Check if required
			isRequired := slices.Contains(parts[1:], "omitempty")

			if isRequired {
				required = append(required, name)
			}

			// Generate schema for the field
			fieldSchema := g.generateSchemaForType(field.Type)
			schema.Properties[name] = fieldSchema
		}

		if len(required) > 0 {
			schema.Required = required
		}

		return schema

	case reflect.Slice, reflect.Array:
		return SchemaObject{
			Type:  "array",
			Items: &SchemaObject{Ref: g.getTypeRef(t.Elem())},
		}

	case reflect.Map:
		return SchemaObject{
			Type: "object",
			// We could add additional properties here if needed
		}

	case reflect.String:
		return SchemaObject{Type: "string"}

	case reflect.Bool:
		return SchemaObject{Type: "boolean"}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return SchemaObject{Type: "integer"}

	case reflect.Float32, reflect.Float64:
		return SchemaObject{Type: "number"}

	default:
		return SchemaObject{Type: "object"}
	}
}

// getTypeRef returns a reference to a component schema for complex types
func (g *OpenAPIGenerator) getTypeRef(t reflect.Type) string {
	// Dereference pointer if needed
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {
		name := t.Name()
		if name != "" {
			// Add schema to components if not already there
			if _, exists := g.schemas[name]; !exists {
				g.schemas[name] = g.generateSchemaForType(t)
			}
			return "#/components/schemas/" + name
		}
	}

	// For simple types, we don't need a reference
	schema := g.generateSchemaForType(t)
	return schema.Type
}

// generateParametersForType converts a struct type to query parameters
func (g *OpenAPIGenerator) generateParametersForType(t reflect.Type) []ParameterObject {
	if t == nil || t.Kind() != reflect.Struct {
		return nil
	}

	var parameters []ParameterObject

	for i := range t.NumField() {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get JSON tag
		tag := field.Tag.Get("json")
		if tag == "-" {
			continue
		}

		// Parse the tag
		parts := strings.Split(tag, ",")
		name := parts[0]
		if name == "" {
			name = field.Name
		}

		// Check if required
		isRequired := slices.Contains(parts[1:], "omitempty")

		// Generate schema for the field
		fieldSchema := g.generateSchemaForType(field.Type)

		// Create parameter
		param := ParameterObject{
			Name:     name,
			In:       "query",
			Required: isRequired,
			Schema:   fieldSchema,
		}

		parameters = append(parameters, param)
	}

	return parameters
}

// GenerateJSON returns the OpenAPI schema as a JSON string
func (g *OpenAPIGenerator) GenerateJSON() (string, error) {
	schema, err := g.GenerateSchema()
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui.css" />
    <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@4.5.0/favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@4.5.0/favicon-16x16.png" sizes="16x16" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
    </style>
</head>

<body>
    <div id="swagger-ui"></div>

    <script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-bundle.js" charset="UTF-8"> </script>
    <script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-standalone-preset.js" charset="UTF-8"> </script>
    <script>
        window.onload = function () {
            const ui = SwaggerUIBundle({
                url: "/swagger.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
            window.ui = ui;
        };
    </script>
</body>
</html>`
