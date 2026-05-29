package docgen

import (
	bytespkg "bytes"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/paleviews/hapi/cmd/protoc-gen-hapi/internal/serviceregistry"
)

func (d *document) yaml() (string, error) {
	var buf bytespkg.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(d.yamlNode()); err != nil {
		_ = enc.Close()
		return "", err
	}
	if err := enc.Close(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (d *document) yamlNode() *yaml.Node {
	pairs := []yamlPair{
		yamlPairForStyledString("openapi", "3.0.3", yaml.DoubleQuotedStyle),
		yamlPairForNode("info", yamlMapping(
			yamlPairForString("version", d.info.version),
			yamlPairForString("title", d.info.title),
		)),
	}
	if len(d.servers) > 0 {
		servers := make([]*yaml.Node, 0, len(d.servers))
		for _, srv := range d.servers {
			servers = append(servers, yamlMapping(yamlPairForString("url", srv.url)))
		}
		pairs = append(pairs, yamlPairForNode("servers", yamlSequence(servers...)))
	}

	paths := make([]yamlPair, 0, len(d.endpoints))
	for _, ep := range d.endpoints {
		paths = append(paths, yamlPairForNode(ep.path, ep.yamlNode()))
	}
	pairs = append(pairs, yamlPairForNode("paths", yamlMapping(paths...)))

	if d.securities.shouldPrint() || d.schemas.shouldPrint() {
		componentPairs := make([]yamlPair, 0, 2)
		if d.securities.shouldPrint() {
			securityPairs := make([]yamlPair, 0, len(d.securities))
			for _, sec := range d.securities {
				securityPairs = append(securityPairs, yamlPairForNode(sec.name, sec.yamlNode()))
			}
			componentPairs = append(componentPairs, yamlPairForNode("securitySchemes", yamlMapping(securityPairs...)))
		}
		if d.schemas.shouldPrint() {
			schemaPairs := make([]yamlPair, 0, len(d.schemas))
			for _, scm := range d.schemas {
				if scm.asAnonymousObject() {
					continue
				}
				schemaPairs = append(schemaPairs, yamlPairForNode(scm.id, scm.componentYAMLNode()))
			}
			componentPairs = append(componentPairs, yamlPairForNode("schemas", yamlMapping(schemaPairs...)))
		}
		pairs = append(pairs, yamlPairForNode("components", yamlMapping(componentPairs...)))
	}

	return yamlMapping(pairs...)
}

func (e *endpoint) yamlNode() *yaml.Node {
	pairs := make([]yamlPair, 0, len(e.operations))
	for _, op := range e.operations {
		pairs = append(pairs, yamlPairForNode(httpVerb(op.httpMethod), op.yamlNode()))
	}
	return yamlMapping(pairs...)
}

func (o *operation) yamlNode() *yaml.Node {
	pairs := make([]yamlPair, 0, 6)
	if desc := o.description.yamlNode(); desc != nil {
		pairs = append(pairs, yamlPairForNode("description", desc))
	}
	pairs = append(pairs, yamlPairForString("operationId", o.operationID))
	if o.security != nil {
		pairs = append(pairs, yamlPairForNode("security", yamlSequence(
			yamlMapping(yamlPairForNode(o.security.name, yamlSequence())),
		)))
	}
	if len(o.parameters) > 0 {
		parameters := make([]*yaml.Node, 0, len(o.parameters))
		for _, param := range o.parameters {
			parameters = append(parameters, parameterYAMLNode(param))
		}
		pairs = append(pairs, yamlPairForNode("parameters", yamlSequence(parameters...)))
	}
	if o.requestBody != nil {
		pairs = append(pairs, yamlPairForNode("requestBody", yamlMapping(
			yamlPairForNode("content", contentYAMLNode(typeYAMLNode(o.requestBody))),
		)))
	}

	responsePairs := []yamlPair{
		yamlPairForString("description", "OK"),
	}
	if len(o.responseHeaders) > 0 {
		headers := make([]yamlPair, 0, len(o.responseHeaders))
		for _, header := range o.responseHeaders {
			headers = append(headers, yamlPairForNode(header.name, header.yamlNode()))
		}
		responsePairs = append(responsePairs, yamlPairForNode("headers", yamlMapping(headers...)))
	}
	responsePairs = append(responsePairs, yamlPairForNode("content", contentYAMLNode(typeYAMLNode(o.response))))
	pairs = append(pairs, yamlPairForNode("responses", yamlMapping(yamlPairForStyledKeyNode(
		"200",
		yaml.SingleQuotedStyle,
		yamlMapping(responsePairs...),
	))))
	return yamlMapping(pairs...)
}

func parameterYAMLNode(param parameter) *yaml.Node {
	switch p := param.(type) {
	case *pathParameter:
		pairs := []yamlPair{
			yamlPairForString("name", p.name),
		}
		if desc := p.description.yamlNode(); desc != nil {
			pairs = append(pairs, yamlPairForNode("description", desc))
		}
		pairs = append(pairs,
			yamlPairForBool("required", true),
			yamlPairForString("in", "path"),
			yamlPairForNode("schema", typeYAMLNode(p.scalarType)),
		)
		return yamlMapping(pairs...)
	case *queryParameter:
		pairs := []yamlPair{
			yamlPairForString("name", p.name),
		}
		if desc := p.description.yamlNode(); desc != nil {
			pairs = append(pairs, yamlPairForNode("description", desc))
		}
		pairs = append(pairs, yamlPairForString("in", "query"))
		if queryParameterUsesSchema(p) {
			pairs = append(pairs, yamlPairForNode("schema", typeYAMLNode(p._type)))
		} else {
			pairs = append(pairs, yamlPairForNode("content", contentYAMLNode(typeYAMLNode(p._type))))
		}
		return yamlMapping(pairs...)
	default:
		panic(fmt.Errorf("unknown parameter type %T", param))
	}
}

func (h *header) yamlNode() *yaml.Node {
	pairs := make([]yamlPair, 0, 2)
	if desc := h.description.yamlNode(); desc != nil {
		pairs = append(pairs, yamlPairForNode("description", desc))
	}
	scm := typeYAMLNode(h.scalarType)
	h.possibilities.addYAMLNodeEnum(scm)
	pairs = append(pairs, yamlPairForNode("schema", scm))
	return yamlMapping(pairs...)
}

func contentYAMLNode(schema *yaml.Node) *yaml.Node {
	return yamlMapping(yamlPairForNode(jsonMediaType, yamlMapping(yamlPairForNode("schema", schema))))
}

func typeYAMLNode(t _type) *yaml.Node {
	switch typ := t.(type) {
	case builtin:
		return typ.yamlNode()
	case enum:
		return typ.yamlNode()
	case bytes:
		return yamlMapping(
			yamlPairForString("type", "string"),
			yamlPairForString("format", "byte"),
		)
	case anonymousObject:
		return typ.yamlNode()
	case *stringKeyedMap:
		return yamlMapping(
			yamlPairForString("type", "object"),
			yamlPairForNode("additionalProperties", typeYAMLNode(typ.value)),
		)
	case *array:
		return yamlMapping(
			yamlPairForString("type", "array"),
			yamlPairForNode("items", typeYAMLNode(typ.item)),
		)
	case *schema:
		if typ.asAnonymousObject() {
			return typeYAMLNode(typ.obj)
		}
		return yamlMapping(yamlPairForStyledString("$ref", "#/components/schemas/"+typ.id, yaml.SingleQuotedStyle))
	default:
		panic(fmt.Errorf("unknown schema type %T", t))
	}
}

func (b builtin) yamlNode() *yaml.Node {
	switch b {
	case builtinBool:
		return yamlMapping(yamlPairForString("type", "boolean"))
	case builtinInt32, builtinUint32:
		return yamlMapping(
			yamlPairForString("type", "integer"),
			yamlPairForString("format", "int32"),
		)
	case builtinInt64, builtinUint64:
		return yamlMapping(
			yamlPairForString("type", "integer"),
			yamlPairForString("format", "int64"),
		)
	case builtinFloat32:
		return yamlMapping(
			yamlPairForString("type", "number"),
			yamlPairForString("format", "float"),
		)
	case builtinFloat64:
		return yamlMapping(
			yamlPairForString("type", "number"),
			yamlPairForString("format", "double"),
		)
	case builtinString:
		return yamlMapping(yamlPairForString("type", "string"))
	default:
		panic(fmt.Errorf("unknown builtin type %d", b))
	}
}

func (e enum) yamlNode() *yaml.Node {
	items := make([]*yaml.Node, 0, len(e))
	for _, item := range e {
		node := yamlString(item.name)
		switch len(item.description) {
		case 0:
		case 1:
			node.LineComment = item.description[0]
		default:
			node.HeadComment = strings.Join(item.description, "\n")
		}
		items = append(items, node)
	}
	return yamlMapping(
		yamlPairForString("type", "string"),
		yamlPairForNode("enum", yamlSequence(items...)),
	)
}

func (ao anonymousObject) yamlNode() *yaml.Node {
	pairs := []yamlPair{
		yamlPairForString("type", "object"),
	}
	if len(ao) == 0 {
		return yamlMapping(pairs...)
	}
	properties := make([]yamlPair, 0, len(ao))
	for _, prop := range ao {
		properties = append(properties, yamlPairForNode(prop.name, prop.yamlNode()))
	}
	pairs = append(pairs, yamlPairForNode("properties", yamlMapping(properties...)))
	return yamlMapping(pairs...)
}

func (prop *property) yamlNode() *yaml.Node {
	pairs := make([]yamlPair, 0, 4)
	if scm, ok := prop._type.(*schema); !ok || scm.asAnonymousObject() {
		if desc := prop.description.yamlNode(); desc != nil {
			pairs = append(pairs, yamlPairForNode("description", desc))
		}
	}
	scm := typeYAMLNode(prop._type)
	pairs = append(pairs, mappingPairs(scm)...)
	prop.possibilities.addYAMLNodeEnumToPairs(&pairs)
	return yamlMapping(pairs...)
}

func (s *schema) componentYAMLNode() *yaml.Node {
	pairs := make([]yamlPair, 0, 3)
	if desc := s.description.yamlNode(); desc != nil {
		pairs = append(pairs, yamlPairForNode("description", desc))
	}
	pairs = append(pairs, mappingPairs(typeYAMLNode(s.obj))...)
	return yamlMapping(pairs...)
}

func (s *security) yamlNode() *yaml.Node {
	return yamlMapping(
		yamlPairForString("type", s._type),
		yamlPairForString("scheme", s.scheme),
	)
}

func (ps possibilities) addYAMLNodeEnum(node *yaml.Node) {
	if len(ps) == 0 {
		return
	}
	pairs := mappingPairs(node)
	ps.addYAMLNodeEnumToPairs(&pairs)
	*node = *yamlMapping(pairs...)
}

func (ps possibilities) addYAMLNodeEnumToPairs(pairs *[]yamlPair) {
	if len(ps) == 0 {
		return
	}
	items := make([]*yaml.Node, 0, len(ps))
	for _, p := range ps {
		item := possibilityYAMLNode(p.literal)
		if p.comment != "" {
			item.LineComment = p.comment
		}
		items = append(items, item)
	}
	*pairs = append(*pairs, yamlPairForNode("enum", yamlSequence(items...)))
}

func possibilityYAMLNode(literal string) *yaml.Node {
	if strings.HasPrefix(literal, "'") && strings.HasSuffix(literal, "'") && len(literal) >= 2 {
		return yamlStyledString(literal[1:len(literal)-1], yaml.SingleQuotedStyle)
	}
	if _, err := strconv.ParseInt(literal, 10, 64); err == nil {
		return yamlInt(literal)
	}
	return yamlString(literal)
}

func (d description) yamlNode() *yaml.Node {
	switch len(d) {
	case 0:
		return nil
	case 1:
		return yamlString(d[0])
	default:
		return yamlLiteral(d.openAPIString())
	}
}

func httpVerb(m serviceregistry.HTTPMethod) string {
	switch m {
	case serviceregistry.HTTPMethodGet:
		return "get"
	case serviceregistry.HTTPMethodPost:
		return "post"
	case serviceregistry.HTTPMethodPut:
		return "put"
	case serviceregistry.HTTPMethodPatch:
		return "patch"
	case serviceregistry.HTTPMethodDelete:
		return "delete"
	default:
		panic(serviceregistry.ErrUnreachableCode)
	}
}

type yamlPair struct {
	key   *yaml.Node
	value *yaml.Node
}

func yamlMapping(pairs ...yamlPair) *yaml.Node {
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}
	for _, pair := range pairs {
		node.Content = append(node.Content, pair.key, pair.value)
	}
	return node
}

func yamlSequence(items ...*yaml.Node) *yaml.Node {
	return &yaml.Node{
		Kind:    yaml.SequenceNode,
		Tag:     "!!seq",
		Content: items,
	}
}

func yamlPairForString(key, value string) yamlPair {
	return yamlPairForNode(key, yamlString(value))
}

func yamlPairForStyledString(key, value string, style yaml.Style) yamlPair {
	return yamlPairForNode(key, yamlStyledString(value, style))
}

func yamlPairForStyledKeyNode(key string, style yaml.Style, value *yaml.Node) yamlPair {
	return yamlPair{
		key:   yamlStyledString(key, style),
		value: value,
	}
}

func yamlPairForBool(key string, value bool) yamlPair {
	return yamlPairForNode(key, yamlBool(value))
}

func yamlPairForNode(key string, value *yaml.Node) yamlPair {
	return yamlPair{
		key:   yamlString(key),
		value: value,
	}
}

func yamlString(value string) *yaml.Node {
	return yamlStyledString(value, 0)
}

func yamlStyledString(value string, style yaml.Style) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: value,
		Style: style,
	}
}

func yamlLiteral(value string) *yaml.Node {
	return yamlStyledString(value, yaml.LiteralStyle)
}

func yamlBool(value bool) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!bool",
		Value: strconv.FormatBool(value),
	}
}

func yamlInt(value string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!int",
		Value: value,
	}
}

func mappingPairs(node *yaml.Node) []yamlPair {
	if node.Kind != yaml.MappingNode {
		panic(fmt.Errorf("expected YAML mapping node, got %v", node.Kind))
	}
	pairs := make([]yamlPair, 0, len(node.Content)/2)
	for i := 0; i < len(node.Content); i += 2 {
		pairs = append(pairs, yamlPair{
			key:   node.Content[i],
			value: node.Content[i+1],
		})
	}
	return pairs
}
