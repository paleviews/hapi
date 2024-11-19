package docgen

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/paleviews/hapi/cmd/protoc-gen-hapi/internal/serviceregistry"
	"github.com/paleviews/hapi/descriptor/annotations"
)

// errInconsistent serves as an input to panic when some states conflict with each other.
// It's better to panic than to generate wrong docs which in most cases are harder to detect.
var errInconsistent = errors.New("inconsistent due to code bugs")

func transform(reg *serviceregistry.Registry, plugin *protogen.Plugin) (*document, error) {
	tf := newTransformer(reg, plugin)
	err := tf.transform()
	if err != nil {
		return nil, err
	}
	var securities securities
	if tf.security != nil {
		securities = append(securities, tf.security)
	}
	return &document{
		info:       tf.docInfo,
		servers:    tf.servers,
		endpoints:  tf.endpoints,
		securities: securities,
		schemas:    tf.schemas,
	}, nil
}

type transformer struct {
	reg    *serviceregistry.Registry
	plugin *protogen.Plugin

	endpoints         []*endpoint
	endpointByRawPath map[string]*endpoint

	schemas    []*schema
	schemaByID map[string]*schema

	docInfo docInfo
	servers []server

	enumValueByResponseCode map[int32]*protogen.EnumValue // lazily populate
	security                *security                     // lazily populate
}

func newTransformer(reg *serviceregistry.Registry, plugin *protogen.Plugin) *transformer {
	return &transformer{
		reg:               reg,
		plugin:            plugin,
		endpointByRawPath: make(map[string]*endpoint),
		schemaByID:        make(map[string]*schema),
	}
}

func (tf *transformer) transform() error {
	if err := tf.checkRouteCollision(); err != nil {
		return err
	}
	tf.docInfo = docInfo{
		version: tf.reg.GlobalInfo().Info.Version,
		title:   tf.reg.GlobalInfo().Info.Title,
	}
	for _, v := range tf.reg.GlobalInfo().Servers {
		tf.servers = append(tf.servers, server{
			url: v.Url,
		})
	}
	for _, rpc := range tf.reg.RPCs() {
		path, op, err := tf.transformRPC(rpc)
		if err != nil {
			return err
		}
		if ep, ok := tf.endpointByRawPath[path]; ok {
			ep.operations = append(ep.operations, op)
			continue
		}
		ep := &endpoint{
			path:       path,
			operations: []*operation{op},
		}
		tf.endpoints = append(tf.endpoints, ep)
		tf.endpointByRawPath[path] = ep
	}
	return nil
}

func (tf *transformer) checkRouteCollision() error {
	type routeKey struct {
		patternPath string
		httpMethod  serviceregistry.HTTPMethod
	}
	routeRecord := make(map[routeKey]*protogen.Method)
	for _, rpc := range tf.reg.RPCs() {
		info := tf.reg.RPCInfo(rpc)
		{
			_routeKey := routeKey{
				patternPath: info.Path.Pattern,
				httpMethod:  info.HTTPMethod,
			}
			if r, ok := routeRecord[_routeKey]; ok {
				return fmt.Errorf("%s at %s shares the same route with %s at %s",
					rpc.Desc.FullName(), rpc.Location.SourceFile,
					r.Desc.FullName(), r.Location.SourceFile)
			}
			routeRecord[_routeKey] = rpc
		}
	}
	return nil
}

func (tf *transformer) transformRPC(rpc *protogen.Method) (path string, _ *operation, _ error) {
	info := tf.reg.RPCInfo(rpc)
	parameters, requestBody, err := tf.transformRPCInput(rpc.Input, info)
	if err != nil {
		return "", nil, err
	}
	response, responseHeaders, err := tf.transformRPCOutput(rpc, info)
	if err != nil {
		return "", nil, err
	}
	var sec *security
	if !info.SkipAuth && tf.reg.GlobalInfo().AuthKind != annotations.AuthKind_AUTH_KIND_NONE {
		if tf.security == nil {
			tf.security = &security{
				name:   "bearerInHeader",
				_type:  "http",
				scheme: "bearer",
			}
		}
		sec = tf.security
	}
	return info.Path.Raw, &operation{
		httpMethod:      info.HTTPMethod,
		description:     tf.parseDescription(string(rpc.Comments.Leading)),
		operationID:     string(rpc.Desc.FullName()),
		security:        sec,
		parameters:      parameters,
		requestBody:     requestBody,
		responseHeaders: responseHeaders,
		response:        response,
		skipAuth:        info.SkipAuth,
		group:           string(rpc.Parent.Desc.FullName()),
	}, nil
}

func (tf *transformer) transformRPCInput(
	input *protogen.Message, rpcInfo *serviceregistry.RPCInfo,
) ([]parameter, object, error) {

	var parameters []parameter
	pathParams := tf.transformPathParams(rpcInfo.PathParams)
	for _, v := range pathParams {
		parameters = append(parameters, v)
	}
	if rpcInfo.AllInputFieldsInPath {
		return parameters, nil, nil
	}

	hasRequestBody := rpcInfo.HTTPMethod != serviceregistry.HTTPMethodGet &&
		rpcInfo.HTTPMethod != serviceregistry.HTTPMethodDelete

	var properties anonymousObject
	if len(rpcInfo.PathParams) == 0 {
		scm, err := tf.transformMessage(input, nil)
		if err != nil {
			return nil, nil, err
		}
		if hasRequestBody {
			return parameters, scm, nil
		}
		properties = scm.obj
	} else {
		exclude := make(map[string]struct{}, len(rpcInfo.PathParams))
		for _, v := range rpcInfo.PathParams {
			exclude[v.PathParam] = struct{}{}
		}
		for _, v := range input.Fields {
			if _, ok := exclude[string(v.Desc.Name())]; ok {
				continue
			}
			p, err := tf.transformField(v, nil)
			if err != nil {
				return nil, nil, err
			}
			properties = append(properties, p)
		}
		if hasRequestBody {
			return parameters, properties, nil
		}
	}

	for _, prop := range properties {
		parameters = append(parameters, &queryParameter{
			name:        prop.name,
			description: prop.description,
			_type:       prop._type,
		})
	}
	return parameters, nil, nil
}

func (tf *transformer) transformRPCOutput(
	rpc *protogen.Method, info *serviceregistry.RPCInfo,
) (object, []*header, error) {

	respScm, err := tf.transformMessage(rpc.Output, nil)
	if err != nil {
		return nil, nil, err
	}

	var (
		codePossibilities possibilities
		descPossibilities possibilities
		codeDesc          description
	)
	if codeInfo := tf.reg.ResponseCodeInfo(); codeInfo != nil && len(info.ResponseCodes) > 0 {
		if tf.enumValueByResponseCode == nil {
			tf.enumValueByResponseCode = make(map[int32]*protogen.EnumValue, len(codeInfo.DescriptionByEnumValue))
			for k := range codeInfo.DescriptionByEnumValue {
				tf.enumValueByResponseCode[int32(k.Desc.Number())] = k
			}
		}
		codeDesc = append(codeDesc, "response code enums:")
		for _, v := range info.ResponseCodes {
			enumValue := tf.enumValueByResponseCode[v]
			if enumValue == nil {
				return nil, nil, fmt.Errorf("unknown response code %d for %s at %s",
					v, rpc.Desc.FullName(), rpc.Location.SourceFile)
			}
			valueDesc := codeInfo.DescriptionByEnumValue[enumValue]
			codeDesc = append(codeDesc, fmt.Sprintf(".. %d: %s", v, valueDesc))
			for _, vv := range tf.parseDescription(string(enumValue.Comments.Leading)) {
				codeDesc = append(codeDesc, fmt.Sprintf(".... %s", vv))
			}
			codePossibilities = append(codePossibilities, &possibility{
				literal: fmt.Sprintf("%d", v),
				comment: valueDesc,
			})
			descPossibilities = append(descPossibilities, &possibility{
				literal: fmt.Sprintf("'%s'", valueDesc),
				comment: fmt.Sprintf("%d", v),
			})
		}
	}

	switch tf.reg.GlobalInfo().ResponseCodeIn {
	case annotations.ResponseCodeLocation_RESPONSE_CODE_LOCATION_BODY:
		return anonymousObject{
			{
				name:          "code",
				_type:         builtinInt32,
				description:   codeDesc,
				possibilities: codePossibilities,
			},
			{
				name:          "message",
				_type:         builtinString,
				possibilities: descPossibilities,
			},
			{
				name:  "data",
				_type: respScm,
			},
		}, nil, nil
	case annotations.ResponseCodeLocation_RESPONSE_CODE_LOCATION_HEADER:
		return respScm, []*header{
			{
				name:          "X-Hapi-Code",
				scalarType:    builtinInt32,
				description:   codeDesc,
				possibilities: codePossibilities,
			},
			{
				name:          "X-Hapi-Message",
				scalarType:    builtinString,
				possibilities: descPossibilities,
			},
		}, nil
	}
	panic(serviceregistry.ErrUnreachableCode)
}

func (tf *transformer) transformMessage(msg *protogen.Message, rb *referencedBy) (*schema, error) {
	if msg.Desc.IsMapEntry() {
		panic(errInconsistent)
	}
	schemaID := string(msg.Desc.FullName())

	if s, ok := tf.schemaByID[schemaID]; ok {
		if s == nil {
			panic(errInconsistent)
		}
		if s.refCount > 0 {
			s.refCount++
		} else {
			ids := rb.findCycle(schemaID)
			switch n := len(ids); n {
			case 0:
				panic(errInconsistent)
			case 1:
			default:
				for _, id := range ids[:n-1] {
					s := tf.schemaByID[id]
					if s == nil {
						panic(errInconsistent)
					}
					s.refCount += 2
				}
			}
			s.refCount += 2
		}
		return s, nil
	}
	scm := &schema{
		id:          schemaID,
		description: tf.parseDescription(string(msg.Comments.Leading)),
	}
	tf.schemaByID[schemaID] = scm
	nextRefBy := rb.add(schemaID)
	properties := make(anonymousObject, 0, len(msg.Fields))
	for _, v := range msg.Fields {
		p, err := tf.transformField(v, nextRefBy)
		if err != nil {
			return nil, err
		}
		properties = append(properties, p)
	}
	scm.refCount++
	scm.obj = properties
	tf.schemas = append(tf.schemas, scm)
	return scm, nil
}

func (tf *transformer) transformField(field *protogen.Field, rb *referencedBy) (*property, error) {
	if field.Desc.IsMap() {
		if field.Message.Fields[0].Desc.Kind() != protoreflect.StringKind {
			return nil, fmt.Errorf("type of map key must be string for field %s at %s",
				field.Desc.FullName(), field.Location.SourceFile)
		}
		valueProperty, err := tf.transformField(field.Message.Fields[1], rb)
		if err != nil {
			return nil, err
		}
		return &property{
			name:        string(field.Desc.Name()),
			description: tf.parseDescription(string(field.Comments.Leading)),
			_type: &stringKeyedMap{
				value: valueProperty._type,
			},
		}, nil
	}
	desc := tf.parseDescription(string(field.Comments.Leading))
	var typ _type
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		typ = builtinBool
	case protoreflect.EnumKind:
		typ = tf.transformEnum(field.Enum)
		desc = tf.meshEnumParentDescription(desc, field.Enum)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		typ = builtinInt32
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		typ = builtinUint32
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		typ = builtinInt64
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		typ = builtinUint64
	case protoreflect.FloatKind:
		typ = builtinFloat32
	case protoreflect.DoubleKind:
		typ = builtinFloat64
	case protoreflect.StringKind:
		typ = builtinString
	case protoreflect.BytesKind:
		typ = bytes{}
	case protoreflect.MessageKind, protoreflect.GroupKind:
		scm, err := tf.transformMessage(field.Message, rb)
		if err != nil {
			return nil, err
		}
		typ = scm
	}
	if field.Desc.IsList() {
		typ = &array{
			item: typ,
		}
	}
	return &property{
		name:        string(field.Desc.Name()),
		description: desc,
		_type:       typ,
	}, nil
}

func (tf *transformer) transformPathParams(infos []*serviceregistry.PathParamInfo) []*pathParameter {
	r := make([]*pathParameter, 0, len(infos))
	for _, v := range infos {
		desc := tf.parseDescription(v.RawComment)
		if v.ScalarType == serviceregistry.ScalarTypeEnum {
			desc = tf.meshEnumParentDescription(desc, v.Enum)
		}
		r = append(r, &pathParameter{
			name:        v.PathParam,
			description: desc,
			scalarType:  tf.transformScalarType(v.ScalarType, v.Enum),
		})
	}
	return r
}

func (tf *transformer) transformScalarType(st serviceregistry.ScalarType, e *protogen.Enum) scalarType {
	switch st {
	case serviceregistry.ScalarTypeBool:
		return builtinBool
	case serviceregistry.ScalarTypeInt32:
		return builtinInt32
	case serviceregistry.ScalarTypeUint32:
		return builtinUint32
	case serviceregistry.ScalarTypeInt64:
		return builtinInt64
	case serviceregistry.ScalarTypeUint64:
		return builtinUint64
	case serviceregistry.ScalarTypeFloat32:
		return builtinFloat32
	case serviceregistry.ScalarTypeFloat64:
		return builtinFloat64
	case serviceregistry.ScalarTypeEnum:
		return tf.transformEnum(e)
	case serviceregistry.ScalarTypeString:
		return builtinString
	default:
		panic(serviceregistry.ErrUnreachableCode)
	}
}

func (tf *transformer) transformEnum(e *protogen.Enum) enum {
	var r enum
	for _, v := range e.Values {
		r = append(r, enumItem{
			name:        string(v.Desc.Name()),
			description: tf.parseDescription(string(v.Comments.Leading)),
			value:       int32(v.Desc.Number()),
		})
	}
	return r
}

func (tf *transformer) meshEnumParentDescription(parent description, enum *protogen.Enum) description {
	if len(parent) > 0 {
		parent = append(parent, "")
	}
	parent = append(parent, fmt.Sprintf(". %s:", enum.Desc.FullName()))
	for _, v := range tf.parseDescription(string(enum.Comments.Leading)) {
		parent = append(parent, fmt.Sprintf("... %s", v))
	}
	parent = append(parent, fmt.Sprintf("... enums:"))
	for _, v := range enum.Values {
		valueDesc := tf.parseDescription(string(v.Comments.Leading))
		switch len(valueDesc) {
		case 0:
			parent = append(parent, fmt.Sprintf("..... %d: %s", v.Desc.Number(), v.Desc.Name()))
		case 1:
			parent = append(parent, fmt.Sprintf("..... %d: %s | %s", v.Desc.Number(), v.Desc.Name(), valueDesc[0]))
		default:
			parent = append(parent, fmt.Sprintf("..... %d: %s", v.Desc.Number(), v.Desc.Name()))
			for _, vv := range valueDesc {
				parent = append(parent, fmt.Sprintf("....... %s", vv))
			}
		}
	}
	return parent
}

func (tf *transformer) parseDescription(comment string) description {
	comment = strings.Trim(comment, "\n\t ")
	if comment == "" {
		return nil
	}
	s := strings.Split(comment, "\n")
	for k, v := range s {
		s[k] = strings.Trim(v, "\t ")
	}
	return s
}
