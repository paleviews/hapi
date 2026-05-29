package docgen

import "github.com/paleviews/hapi/cmd/protoc-gen-hapi/internal/serviceregistry"

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

type endpoint struct {
	path       string
	operations []*operation
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
}

type description []string

type parameter interface {
	isParameter()
}

type pathParameter struct {
	name        string
	description description
	scalarType  scalarType
}

func (*pathParameter) isParameter() {}

type queryParameter struct {
	name        string
	description description
	_type       _type
}

func (*queryParameter) isParameter() {}

type _type interface {
	isType()
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

type property struct {
	name          string
	description   description
	_type         _type
	possibilities possibilities
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

type enumItem struct {
	name        string
	description description
	value       int32
}

type enum []enumItem

func (enum) isType() {}

func (enum) isScalar() {}

type bytes struct{}

func (bytes) isType() {}

type anonymousObject []*property

func (anonymousObject) isType() {}

func (anonymousObject) isObject() {}

type stringKeyedMap struct {
	value _type
}

func (*stringKeyedMap) isType() {}

type array struct {
	item _type
}

func (*array) isType() {}

type header struct {
	name          string
	description   description
	scalarType    scalarType
	possibilities possibilities
}

type schema struct {
	id          string
	description description
	refCount    int64
	obj         anonymousObject
}

func (*schema) isType() {}

func (*schema) isObject() {}

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

type security struct {
	name   string
	_type  string
	scheme string
}

type securities []*security

func (ss securities) shouldPrint() bool {
	return len(ss) > 0
}
