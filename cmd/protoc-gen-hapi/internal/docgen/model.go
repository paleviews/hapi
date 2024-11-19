package docgen

import (
	"fmt"

	"github.com/paleviews/hapi/cmd/protoc-gen-hapi/internal/printer"
	"github.com/paleviews/hapi/cmd/protoc-gen-hapi/internal/serviceregistry"
)

type docInfo struct {
	version string
	title   string
}

type server struct {
	url string
}

type document struct {
	info       docInfo
	servers    []server
	endpoints  []*endpoint
	securities securities
	schemas    schemas
}

func (d *document) print(p *printer.Printer) {
	p.PrintlnWithIndents(`openapi: "3.0.3"`)
	p.PrintlnWithIndents("info:")
	p.IncreaseIndentLevel()
	p.PrintlnWithIndents("version: ", d.info.version)
	p.PrintlnWithIndents("title: ", d.info.title)
	p.DecreaseIndentLevel()
	if len(d.servers) > 0 {
		p.PrintlnWithIndents("servers:")
		p.IncreaseIndentLevel()
		for _, v := range d.servers {
			p.PrintlnWithIndents("- url: ", v.url)
		}
		p.DecreaseIndentLevel()
	}

	p.PrintlnWithIndents("paths:")
	p.IncreaseIndentLevel()
	for _, v := range d.endpoints {
		v.print(p)
	}
	p.DecreaseIndentLevel()

	printSecurities := d.securities.shouldPrint()
	printSchemas := d.schemas.shouldPrint()
	if !printSecurities && !printSchemas {
		return
	}
	p.PrintlnWithIndents("components:")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	if printSecurities {
		d.securities.print(p)
	}
	if printSchemas {
		d.schemas.print(p)
	}
}

type endpoint struct {
	path       string
	operations []*operation
}

func (e *endpoint) print(p *printer.Printer) {
	p.PrintlnWithIndents(e.path, ":")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	for _, v := range e.operations {
		v.print(p)
	}
}

type operation struct {
	httpMethod      serviceregistry.HTTPMethod
	description     description
	operationID     string
	security        *security
	parameters      []parameter
	requestBody     object
	responseHeaders []*header
	response        object
	skipAuth        bool
	group           string
}

func (o *operation) print(p *printer.Printer) {
	getHTTPVerb := func(m serviceregistry.HTTPMethod) string {
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
	p.PrintlnWithIndents(getHTTPVerb(o.httpMethod), ":")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()

	o.description.print(p)

	p.PrintlnWithIndents("operationId: ", o.operationID)

	if o.security != nil {
		p.PrintlnWithIndents("security:")
		p.IncreaseIndentLevel()
		p.PrintlnWithIndents("- ", o.security.name, ": []")
		p.DecreaseIndentLevel()
	}

	if len(o.parameters) > 0 {
		p.PrintlnWithIndents("parameters:")
		p.IncreaseIndentLevel()
		for _, v := range o.parameters {
			v.print(p)
		}
		p.DecreaseIndentLevel()
	}

	if o.requestBody != nil {
		p.PrintlnWithIndents("requestBody:")
		p.IncreaseIndentLevel()
		p.PrintlnWithIndents("content:")
		p.IncreaseIndentLevel()
		p.PrintlnWithIndents("application/json; charset=utf-8:")
		p.IncreaseIndentLevel()
		p.PrintlnWithIndents("schema:")
		p.IncreaseIndentLevel()
		o.requestBody.print(p)
		p.DecreaseIndentLevel()
		p.DecreaseIndentLevel()
		p.DecreaseIndentLevel()
		p.DecreaseIndentLevel()
	}

	p.PrintlnWithIndents("responses:")
	p.IncreaseIndentLevel()
	p.PrintlnWithIndents("'200':")
	p.IncreaseIndentLevel()
	p.PrintlnWithIndents("description: OK")

	if len(o.responseHeaders) > 0 {
		p.PrintlnWithIndents("headers:")
		p.IncreaseIndentLevel()
		for _, v := range o.responseHeaders {
			v.print(p)
		}
		p.DecreaseIndentLevel()
	}

	p.PrintlnWithIndents("content:")
	p.IncreaseIndentLevel()
	p.PrintlnWithIndents("application/json; charset=utf-8:")
	p.IncreaseIndentLevel()
	p.PrintlnWithIndents("schema:")
	p.IncreaseIndentLevel()
	o.response.print(p)
	p.DecreaseIndentLevel()
	p.DecreaseIndentLevel()
	p.DecreaseIndentLevel()
	p.DecreaseIndentLevel()
	p.DecreaseIndentLevel()

	if o.group != "" {
		p.PrintlnWithIndents("tags:")
		p.IncreaseIndentLevel()
		p.PrintlnWithIndents("- ", o.group)
		p.DecreaseIndentLevel()
	}
}

type description []string

func (d description) print(p *printer.Printer) {
	switch n := len(d); n {
	case 0:
	case 1:
		p.PrintlnWithIndents("description: ", d[0])
	default:
		p.PrintlnWithIndents("description: |")
		p.IncreaseIndentLevel()
		defer p.DecreaseIndentLevel()
		for k, v := range d {
			if k == n-1 {
				p.PrintlnWithIndents(v)
			} else {
				p.PrintlnWithIndents(v, `\`)
			}
		}
	}
}

type parameter interface {
	isParameter()
	print(p *printer.Printer)
}

type pathParameter struct {
	name        string
	description description
	scalarType  scalarType
}

func (*pathParameter) isParameter() {}

func (pp *pathParameter) print(p *printer.Printer) {
	p.PrintlnWithIndents("- name: ", pp.name)
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	pp.description.print(p)
	p.PrintlnWithIndents("required: true")
	p.PrintlnWithIndents("in: path")
	p.PrintlnWithIndents("schema:")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	pp.scalarType.print(p)
}

type queryParameter struct {
	name        string
	description description
	_type       _type
}

func (*queryParameter) isParameter() {}

func (qp *queryParameter) print(p *printer.Printer) {
	p.PrintlnWithIndents("- name: ", qp.name)
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	qp.description.print(p)
	p.PrintlnWithIndents("in: query")
	switch qp._type.(type) {
	case scalarType, bytes, *array:
	default:
		p.PrintlnWithIndents("content:")
		p.IncreaseIndentLevel()
		defer p.DecreaseIndentLevel()
		p.PrintlnWithIndents("application/json; charset=utf-8:")
		p.IncreaseIndentLevel()
		defer p.DecreaseIndentLevel()
	}
	p.PrintlnWithIndents("schema:")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	qp._type.print(p)
}

type _type interface {
	isType()
	print(p *printer.Printer)
}

type scalarType interface {
	isScalar()
	_type
}

type object interface {
	isObject()
	_type
}

type possibility struct {
	literal string
	comment string
}

type possibilities []*possibility // type name enum is taken :(

func (ps possibilities) print(p *printer.Printer) {
	if len(ps) == 0 {
		return
	}
	p.PrintlnWithIndents("enum:")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	for _, v := range ps {
		if v.comment == "" {
			p.PrintlnWithIndents("- ", v.literal)
		} else {
			p.PrintlnWithIndents("- ", v.literal, " # ", v.comment)
		}
	}
}

type property struct {
	name          string
	description   description
	_type         _type
	possibilities possibilities
}

func (prop *property) print(p *printer.Printer) {
	p.PrintlnWithIndents(prop.name, ":")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	// $ref cannot have any siblings
	if scm, ok := prop._type.(*schema); !ok || scm.asAnonymousObject() {
		prop.description.print(p)
	}
	prop._type.print(p)
	prop.possibilities.print(p)
}

type builtin uint8

const (
	builtinBool builtin = iota + 1
	builtinInt32
	builtinUint32
	builtinInt64
	builtinUint64
	builtinFloat32
	builtinFloat64
	builtinString
)

func (builtin) isType() {}

func (b builtin) isScalar() {}

func (b builtin) print(p *printer.Printer) {
	switch b {
	case builtinBool:
		p.PrintlnWithIndents("type: boolean")
	case builtinInt32, builtinUint32:
		p.PrintlnWithIndents("type: integer")
		p.PrintlnWithIndents("format: int32")
	case builtinInt64, builtinUint64:
		p.PrintlnWithIndents("type: integer")
		p.PrintlnWithIndents("format: int64")
	case builtinFloat32:
		p.PrintlnWithIndents("type: number")
		p.PrintlnWithIndents("format: float")
	case builtinFloat64:
		p.PrintlnWithIndents("type: number")
		p.PrintlnWithIndents("format: double")
	case builtinString:
		p.PrintlnWithIndents("type: string")
	default:
		panic(fmt.Errorf("unknown builtin type %d", b))
	}
}

type enumItem struct {
	name        string
	description description
	value       int32
}

type enum []enumItem

func (enum) isType() {}

func (enum) isScalar() {}

func (e enum) print(p *printer.Printer) {
	p.PrintlnWithIndents("type: integer")
	p.PrintlnWithIndents("format: int32")
	p.PrintlnWithIndents("enum:")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	for _, v := range e {
		switch len(v.description) {
		case 0:
			p.PrintlnWithIndents(fmt.Sprintf("- %d # %s", v.value, v.name))
		case 1:
			p.PrintlnWithIndents(fmt.Sprintf("- %d # %s: %s", v.value, v.name, v.description[0]))
		default:
			p.PrintlnWithIndents(fmt.Sprintf("# %s:", v.name))
			for _, vv := range v.description {
				p.PrintlnWithIndents(fmt.Sprintf("#   %s", vv))
			}
			p.PrintlnWithIndents(fmt.Sprintf("- %d", v.value))
		}
	}
}

type bytes struct{}

func (bytes) isType() {}

func (bytes) print(p *printer.Printer) {
	p.PrintlnWithIndents("type: string")
	p.PrintlnWithIndents("format: byte")
}

type anonymousObject []*property

func (anonymousObject) isType() {}

func (anonymousObject) isObject() {}

func (ao anonymousObject) print(p *printer.Printer) {
	p.PrintlnWithIndents("type: object")
	if len(ao) == 0 {
		return
	}
	p.PrintlnWithIndents("properties:")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	for _, v := range ao {
		v.print(p)
	}
}

type stringKeyedMap struct {
	value _type
}

func (*stringKeyedMap) isType() {}

func (m *stringKeyedMap) print(p *printer.Printer) {
	p.PrintlnWithIndents("type: object")
	p.PrintlnWithIndents("additionalProperties:")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	m.value.print(p)
}

type array struct {
	item _type
}

func (*array) isType() {}

func (a *array) print(p *printer.Printer) {
	p.PrintlnWithIndents("type: array")
	p.PrintlnWithIndents("items:")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	a.item.print(p)
}

type header struct {
	name          string
	description   description
	scalarType    scalarType
	possibilities possibilities
}

func (h *header) print(p *printer.Printer) {
	p.PrintlnWithIndents(h.name, ":")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	h.description.print(p)
	p.PrintlnWithIndents("schema:")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	h.scalarType.print(p)
	h.possibilities.print(p)
}

type schema struct {
	id          string
	description description
	refCount    int64
	obj         anonymousObject
}

func (*schema) isType() {}

func (*schema) isObject() {}

func (s *schema) print(p *printer.Printer) {
	if s.asAnonymousObject() { // TODO: the description is ignored for now, should combined to the parent description
		s.obj.print(p)
		return
	}
	p.PrintlnWithIndents("$ref: '#/components/schemas/", s.id, "'")
}

func (s *schema) asAnonymousObject() bool {
	return s.refCount == 1
}

type schemas []*schema

func (ss schemas) shouldPrint() bool {
	for _, v := range ss {
		if !v.asAnonymousObject() {
			return true
		}
	}
	return false
}

func (ss schemas) print(p *printer.Printer) {
	p.PrintlnWithIndents("schemas:")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	for _, s := range ss {
		if s.asAnonymousObject() {
			continue
		}
		p.PrintlnWithIndents(s.id, ":")
		p.IncreaseIndentLevel()
		s.description.print(p)
		s.obj.print(p)
		p.DecreaseIndentLevel()
	}
}

type security struct {
	name   string
	_type  string
	scheme string
}

func (s *security) print(p *printer.Printer) {
	p.PrintlnWithIndents(s.name, ":")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	p.PrintlnWithIndents("type: ", s._type)
	p.PrintlnWithIndents("scheme: ", s.scheme)
}

type securities []*security

func (ss securities) shouldPrint() bool {
	return len(ss) > 0
}

func (ss securities) print(p *printer.Printer) {
	p.PrintlnWithIndents("securitySchemes:")
	p.IncreaseIndentLevel()
	defer p.DecreaseIndentLevel()
	for _, s := range ss {
		s.print(p)
	}
}
