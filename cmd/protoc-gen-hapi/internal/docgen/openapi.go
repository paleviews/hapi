package docgen

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/paleviews/hapi/cmd/protoc-gen-hapi/internal/serviceregistry"
)

const jsonMediaType = "application/json; charset=utf-8"

func validateOpenAPI(spec *openapi3.T) error {
	loader := openapi3.NewLoader()
	if err := loader.ResolveRefsIn(spec, nil); err != nil {
		return err
	}
	return spec.Validate(context.Background())
}

func (d *document) openAPI() *openapi3.T {
	spec := &openapi3.T{
		OpenAPI: "3.0.3",
		Info: &openapi3.Info{
			Version: d.info.version,
			Title:   d.info.title,
		},
		Paths: openapi3.NewPaths(),
	}
	for _, srv := range d.servers {
		spec.Servers = append(spec.Servers, &openapi3.Server{URL: srv.url})
	}
	for _, ep := range d.endpoints {
		spec.Paths.Set(ep.path, ep.openAPI())
	}

	if d.securities.shouldPrint() || d.schemas.shouldPrint() {
		components := openapi3.NewComponents()
		if d.securities.shouldPrint() {
			components.SecuritySchemes = make(openapi3.SecuritySchemes, len(d.securities))
			for _, sec := range d.securities {
				components.SecuritySchemes[sec.name] = &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type:   sec._type,
						Scheme: sec.scheme,
					},
				}
			}
		}
		if d.schemas.shouldPrint() {
			components.Schemas = make(openapi3.Schemas)
			for _, scm := range d.schemas {
				if scm.asAnonymousObject() {
					continue
				}
				components.Schemas[scm.id] = scm.componentSchemaRef()
			}
		}
		spec.Components = &components
	}
	return spec
}

func (e *endpoint) openAPI() *openapi3.PathItem {
	item := &openapi3.PathItem{}
	for _, op := range e.operations {
		switch op.httpMethod {
		case serviceregistry.HTTPMethodGet:
			item.Get = op.openAPI()
		case serviceregistry.HTTPMethodPost:
			item.Post = op.openAPI()
		case serviceregistry.HTTPMethodPut:
			item.Put = op.openAPI()
		case serviceregistry.HTTPMethodPatch:
			item.Patch = op.openAPI()
		case serviceregistry.HTTPMethodDelete:
			item.Delete = op.openAPI()
		default:
			panic(serviceregistry.ErrUnreachableCode)
		}
	}
	return item
}

func (o *operation) openAPI() *openapi3.Operation {
	op := &openapi3.Operation{
		Description: o.description.openAPIString(),
		OperationID: o.operationID,
		Responses:   openapi3.NewResponses(),
	}
	if o.security != nil {
		requirements := openapi3.NewSecurityRequirements()
		requirements.With(openapi3.NewSecurityRequirement().Authenticate(o.security.name))
		op.Security = requirements
	}
	if len(o.parameters) > 0 {
		op.Parameters = make(openapi3.Parameters, 0, len(o.parameters))
		for _, param := range o.parameters {
			op.Parameters = append(op.Parameters, openAPIParameter(param))
		}
	}
	if o.requestBody != nil {
		op.RequestBody = &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Content: openAPIContent(openAPISchemaRef(o.requestBody)),
			},
		}
	}

	response := &openapi3.Response{
		Description: openapi3.Ptr("OK"),
		Content:     openAPIContent(openAPISchemaRef(o.response)),
	}
	if len(o.responseHeaders) > 0 {
		response.Headers = make(openapi3.Headers, len(o.responseHeaders))
		for _, header := range o.responseHeaders {
			response.Headers[header.name] = header.openAPI()
		}
	}
	op.Responses.Set("200", &openapi3.ResponseRef{Value: response})
	return op
}

func openAPIParameter(param parameter) *openapi3.ParameterRef {
	switch p := param.(type) {
	case *pathParameter:
		return &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name:        p.name,
				Description: p.description.openAPIString(),
				Required:    true,
				In:          openapi3.ParameterInPath,
				Schema:      openAPISchemaRef(p.scalarType),
			},
		}
	case *queryParameter:
		value := &openapi3.Parameter{
			Name:        p.name,
			Description: p.description.openAPIString(),
			In:          openapi3.ParameterInQuery,
		}
		if queryParameterUsesSchema(p) {
			value.Schema = openAPISchemaRef(p._type)
		} else {
			value.Content = openAPIContent(openAPISchemaRef(p._type))
		}
		return &openapi3.ParameterRef{Value: value}
	default:
		panic(fmt.Errorf("unknown parameter type %T", param))
	}
}

func queryParameterUsesSchema(qp *queryParameter) bool {
	switch qp._type.(type) {
	case scalarType, bytes, *array:
		return true
	default:
		return false
	}
}

func (h *header) openAPI() *openapi3.HeaderRef {
	ref := openAPISchemaRef(h.scalarType)
	h.possibilities.applyOpenAPIEnum(ref)
	return &openapi3.HeaderRef{
		Value: &openapi3.Header{
			Parameter: openapi3.Parameter{
				Description: h.description.openAPIString(),
				Schema:      ref,
			},
		},
	}
}

func openAPIContent(schemaRef *openapi3.SchemaRef) openapi3.Content {
	return openapi3.Content{
		jsonMediaType: &openapi3.MediaType{
			Schema: schemaRef,
		},
	}
}

func openAPISchemaRef(t _type) *openapi3.SchemaRef {
	switch typ := t.(type) {
	case builtin:
		return &openapi3.SchemaRef{Value: typ.openAPISchema()}
	case enum:
		return &openapi3.SchemaRef{Value: typ.openAPISchema()}
	case bytes:
		return &openapi3.SchemaRef{Value: openapi3.NewBytesSchema()}
	case anonymousObject:
		return typ.openAPISchemaRef()
	case *stringKeyedMap:
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeObject},
				AdditionalProperties: openapi3.AdditionalProperties{
					Schema: openAPISchemaRef(typ.value),
				},
			},
		}
	case *array:
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:  &openapi3.Types{openapi3.TypeArray},
				Items: openAPISchemaRef(typ.item),
			},
		}
	case *schema:
		if typ.asAnonymousObject() {
			return openAPISchemaRef(typ.obj)
		}
		return &openapi3.SchemaRef{Ref: "#/components/schemas/" + typ.id}
	default:
		panic(fmt.Errorf("unknown schema type %T", t))
	}
}

func (b builtin) openAPISchema() *openapi3.Schema {
	switch b {
	case builtinBool:
		return openapi3.NewBoolSchema()
	case builtinInt32, builtinUint32:
		return openapi3.NewInt32Schema()
	case builtinInt64, builtinUint64:
		return openapi3.NewInt64Schema()
	case builtinFloat32:
		return (&openapi3.Schema{Type: &openapi3.Types{openapi3.TypeNumber}}).WithFormat("float")
	case builtinFloat64:
		return openapi3.NewFloat64Schema().WithFormat("double")
	case builtinString:
		return openapi3.NewStringSchema()
	default:
		panic(fmt.Errorf("unknown builtin type %d", b))
	}
}

func (e enum) openAPISchema() *openapi3.Schema {
	values := make([]any, 0, len(e))
	for _, item := range e {
		values = append(values, item.name)
	}
	return openapi3.NewStringSchema().WithEnum(values...)
}

func (ao anonymousObject) openAPISchemaRef() *openapi3.SchemaRef {
	scm := &openapi3.Schema{
		Type:       &openapi3.Types{openapi3.TypeObject},
		Properties: make(openapi3.Schemas, len(ao)),
	}
	for _, prop := range ao {
		ref := openAPISchemaRef(prop._type)
		if ref.Ref == "" && ref.Value != nil {
			ref.Value.Description = prop.description.openAPIString()
		}
		prop.possibilities.applyOpenAPIEnum(ref)
		scm.Properties[prop.name] = ref
	}
	return &openapi3.SchemaRef{Value: scm}
}

func (s *schema) componentSchemaRef() *openapi3.SchemaRef {
	ref := openAPISchemaRef(s.obj)
	if ref.Value != nil {
		ref.Value.Description = s.description.openAPIString()
	}
	return ref
}

func (ps possibilities) applyOpenAPIEnum(ref *openapi3.SchemaRef) {
	if len(ps) == 0 || ref == nil || ref.Value == nil {
		return
	}
	ref.Value.Enum = make([]any, 0, len(ps))
	for _, p := range ps {
		ref.Value.Enum = append(ref.Value.Enum, openAPIEnumLiteral(p.literal))
	}
}

func openAPIEnumLiteral(literal string) any {
	if strings.HasPrefix(literal, "'") && strings.HasSuffix(literal, "'") && len(literal) >= 2 {
		return literal[1 : len(literal)-1]
	}
	if n, err := strconv.ParseInt(literal, 10, 64); err == nil {
		return n
	}
	return literal
}

func (d description) openAPIString() string {
	switch len(d) {
	case 0:
		return ""
	case 1:
		return d[0]
	default:
		var b strings.Builder
		for i, line := range d {
			b.WriteString(line)
			if i != len(d)-1 {
				b.WriteByte('\\')
			}
			b.WriteByte('\n')
		}
		return b.String()
	}
}
