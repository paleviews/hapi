package serviceregistry

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/paleviews/hapi/descriptor/annotations"
)

var ErrUnreachableCode = errors.New("unreachable code because the scenario is excluded by the registry")

type ScalarType uint8

const (
	ScalarTypeBool ScalarType = iota + 1
	ScalarTypeInt32
	ScalarTypeUint32
	ScalarTypeInt64
	ScalarTypeUint64
	ScalarTypeFloat32
	ScalarTypeFloat64
	ScalarTypeEnum
	ScalarTypeString
)

type HTTPMethod uint8

const (
	HTTPMethodGet HTTPMethod = iota + 1
	HTTPMethodPost
	HTTPMethodPut
	HTTPMethodPatch
	HTTPMethodDelete
)

type PathParamInfo struct {
	GoField    string
	PathParam  string // also the protobuf field name
	ScalarType ScalarType
	RawComment string
	Enum       *protogen.Enum
}

type PathInfo struct {
	Raw     string
	Pattern string
}

type RPCInfo struct {
	HTTPMethod           HTTPMethod
	Path                 PathInfo
	SkipAuth             bool
	PathParams           []*PathParamInfo
	AllInputFieldsInPath bool
	ResponseCodes        []int32
}

type ResponseCodeInfo struct {
	OriginEnum             *protogen.Enum
	ValueOK                *protogen.EnumValue
	ValueInvalidInput      *protogen.EnumValue
	ValueUnauthenticated   *protogen.EnumValue
	ValueServerError       *protogen.EnumValue
	DescriptionByEnumValue map[*protogen.EnumValue]string
}

type Registry struct {
	extensionTypes   *protoregistry.Types
	rpcList          []*protogen.Method
	rpcInfos         map[*protogen.Method]*RPCInfo
	services         []*protogen.Service
	globalInfo       *annotations.GlobalOptions
	responseCodeInfo *ResponseCodeInfo

	responseCodesMethodExtensionFieldNo int32
}

func Load(plugin *protogen.Plugin) (*Registry, error) {
	reg := &Registry{
		extensionTypes: new(protoregistry.Types),
		rpcInfos:       make(map[*protogen.Method]*RPCInfo),
	}
	err := reg.load(plugin)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

func (reg *Registry) RPCs() []*protogen.Method {
	return reg.rpcList
}

func (reg *Registry) RPCInfo(method *protogen.Method) *RPCInfo {
	i := reg.rpcInfos[method]
	if i == nil {
		panic(ErrUnreachableCode)
	}
	return i
}

func (reg *Registry) GlobalInfo() *annotations.GlobalOptions {
	return reg.globalInfo
}

func (reg *Registry) ResponseCodeInfo() *ResponseCodeInfo {
	return reg.responseCodeInfo
}

func (reg *Registry) load(plugin *protogen.Plugin) error {
	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}
		err := reg.loadGlobalOptions(file)
		if err != nil {
			return err
		}
		for _, v := range file.Enums {
			err := reg.loadResponseCode(v)
			if err != nil {
				return err
			}
		}
		for _, extension := range file.Extensions {
			err := reg.loadExtension(extension)
			if err != nil {
				return err
			}
		}
	}
	if reg.globalInfo == nil {
		return errors.New("global options not specified")
	}
	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}
		for _, service := range file.Services {
			err := reg.loadService(service)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (reg *Registry) loadResponseCode(enum *protogen.Enum) error {
	isResponseCode := proto.GetExtension(enum.Desc.Options(), annotations.E_IsTypeOfResponseCode).(bool)
	if !isResponseCode {
		return nil
	}
	if reg.responseCodeInfo != nil {
		return fmt.Errorf("duplicated response code annotated for %s at %s",
			enum.Desc.FullName(), enum.Location.SourceFile)
	}
	info := ResponseCodeInfo{
		OriginEnum:             enum,
		DescriptionByEnumValue: make(map[*protogen.EnumValue]string, len(enum.Values)),
	}
	for _, v := range enum.Values {
		options := proto.GetExtension(v.Desc.Options(), annotations.E_ResponseCodeValue).(*annotations.ResponseCodeValueOptions)
		if options == nil {
			return fmt.Errorf("no response code options for response code value %s at %s",
				v.Desc.FullName(), v.Location.SourceFile)
		}
		if options.Desc == "" {
			return fmt.Errorf("no description for response code value %s at %s",
				v.Desc.FullName(), v.Location.SourceFile)
		}
		if strings.Contains(options.Desc, "\n") {
			return fmt.Errorf("new line in description for response code value %s at %s",
				v.Desc.FullName(), v.Location.SourceFile)
		}
		info.DescriptionByEnumValue[v] = options.Desc
		switch {
		case options.GetIsOk():
			if info.ValueOK != nil {
				return fmt.Errorf("duplicated response code value 'ok' annotated for %s at %s",
					v.Desc.FullName(), v.Location.SourceFile)
			}
			info.ValueOK = v
		case options.GetIsInvalidInput():
			if info.ValueInvalidInput != nil {
				return fmt.Errorf("duplicated response code value 'invalid_input' annotated for %s at %s",
					v.Desc.FullName(), v.Location.SourceFile)
			}
			info.ValueInvalidInput = v
		case options.GetIsUnauthenticated():
			if info.ValueUnauthenticated != nil {
				return fmt.Errorf("duplicated response code value 'unauthenticated' annotated for %s at %s",
					v.Desc.FullName(), v.Location.SourceFile)
			}
			info.ValueUnauthenticated = v
		case options.GetIsServerError():
			if info.ValueServerError != nil {
				return fmt.Errorf("duplicated response code value 'server_error' annotated for %s at %s",
					v.Desc.FullName(), v.Location.SourceFile)
			}
			info.ValueServerError = v
		}
	}
	switch {
	case info.ValueOK == nil:
		return fmt.Errorf("no predefined response code value 'ok' annotated for %s at %s",
			enum.Desc.FullName(), enum.Location.SourceFile)
	case info.ValueInvalidInput == nil:
		return fmt.Errorf("no predefined response code value 'invalid_input' annotated for %s at %s",
			enum.Desc.FullName(), enum.Location.SourceFile)
	case info.ValueUnauthenticated == nil:
		return fmt.Errorf("no predefined response code value 'unauthenticated' annotated for %s at %s",
			enum.Desc.FullName(), enum.Location.SourceFile)
	case info.ValueServerError == nil:
		return fmt.Errorf("no predefined response code value 'server_error' annotated for %s at %s",
			enum.Desc.FullName(), enum.Location.SourceFile)
	}
	reg.responseCodeInfo = &info
	return nil
}

func (reg *Registry) loadGlobalOptions(file *protogen.File) error {
	options := proto.GetExtension(file.Desc.Options(), annotations.E_GlobalOptions).(*annotations.GlobalOptions)
	if options == nil {
		return nil
	}
	if reg.globalInfo != nil {
		return fmt.Errorf("duplicated global options in %s", file.Desc.Path())
	}
	switch options.ResponseCodeIn {
	case annotations.ResponseCodeLocation_RESPONSE_CODE_LOCATION_BODY:
	case annotations.ResponseCodeLocation_RESPONSE_CODE_LOCATION_HEADER:
	default:
		return fmt.Errorf("unsupported response_code_in %d for global options at %s",
			options.ResponseCodeIn, file.Desc.Path())
	}
	switch options.AuthKind {
	case annotations.AuthKind_AUTH_KIND_NONE:
	case annotations.AuthKind_AUTH_KIND_BEARER_IN_HEADER:
	default:
		return fmt.Errorf("unsupported auth_kind %d for global options at %s",
			options.AuthKind, file.Desc.Path())
	}
	if options.Info == nil {
		return fmt.Errorf("no info in global options at %s", file.Desc.Path())
	}
	if options.Info.Version == "" {
		return fmt.Errorf("no version info in global options at %s", file.Desc.Path())
	}
	if options.Info.Title == "" {
		return fmt.Errorf("no title info in global options at %s", file.Desc.Path())
	}
	reg.globalInfo = options
	return nil
}

func (reg *Registry) loadExtension(ext *protogen.Extension) error {
	if !proto.GetExtension(ext.Desc.Options(), annotations.E_MethodExtensionResponseCodes).(bool) {
		return nil
	}
	if reg.responseCodesMethodExtensionFieldNo > 0 {
		return fmt.Errorf("duplicated method extension for response codes annotated for %s at %s",
			ext.Desc.FullName(), ext.Location.SourceFile)
	}
	if !ext.Desc.IsList() {
		return fmt.Errorf("method extension field (%s at %s) for response codes is not a list",
			ext.Desc.FullName(), ext.Location.SourceFile)
	}
	if ext.Desc.Kind() != protoreflect.EnumKind {
		return fmt.Errorf("method extension field (%s at %s) for response codes is not a list of enums",
			ext.Desc.FullName(), ext.Location.SourceFile)
	}
	if ext.Desc.Enum() != reg.responseCodeInfo.OriginEnum.Desc {
		return fmt.Errorf(
			"type of method extension field (%s at %s) for response codes "+
				"doesn't match the enum annotated as response code %s",
			ext.Desc.FullName(), ext.Location.SourceFile,
			reg.responseCodeInfo.OriginEnum.Desc.FullName())
	}
	reg.responseCodesMethodExtensionFieldNo = int32(ext.Desc.Number())
	err := reg.extensionTypes.RegisterExtension(dynamicpb.NewExtensionType(ext.Desc))
	if err != nil {
		return err
	}
	return nil
}

func (reg *Registry) loadService(service *protogen.Service) error {
	for _, rpc := range service.Methods {
		err := reg.loadRPC(rpc)
		if err != nil {
			return err
		}
	}
	reg.services = append(reg.services, service)
	return nil
}

// https://github.com/golang/protobuf/issues/1260
func (reg *Registry) getResponseCodes(method *protogen.Method) []int32 {
	options := method.Desc.Options().(*descriptorpb.MethodOptions)
	b, err := proto.Marshal(options)
	if err != nil {
		panic(err)
	}
	options = new(descriptorpb.MethodOptions)
	err = proto.UnmarshalOptions{Resolver: reg.extensionTypes}.Unmarshal(b, options)
	if err != nil {
		panic(err)
	}
	var list []int32
	options.ProtoReflect().Range(func(descriptor protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		if !descriptor.IsExtension() {
			return true
		}
		if int32(descriptor.Number()) != reg.responseCodesMethodExtensionFieldNo {
			return true
		}
		for i := 0; i < value.List().Len(); i++ {
			v := value.List().Get(i)
			list = append(list, int32(v.Enum()))
		}
		return false
	})
	return list
}

func (reg *Registry) loadRPC(method *protogen.Method) error {
	options := proto.GetExtension(method.Desc.Options(), annotations.E_Method).(*annotations.MethodOptions)
	info := RPCInfo{
		SkipAuth:      options.SkipAuth,
		ResponseCodes: reg.getResponseCodes(method),
	}
	switch {
	case options.GetGet() != "":
		info.HTTPMethod = HTTPMethodGet
		info.Path.Raw = options.GetGet()
	case options.GetPost() != "":
		info.HTTPMethod = HTTPMethodPost
		info.Path.Raw = options.GetPost()
	case options.GetPut() != "":
		info.HTTPMethod = HTTPMethodPut
		info.Path.Raw = options.GetPut()
	case options.GetPatch() != "":
		info.HTTPMethod = HTTPMethodPatch
		info.Path.Raw = options.GetPatch()
	case options.GetDelete() != "":
		info.HTTPMethod = HTTPMethodDelete
		info.Path.Raw = options.GetDelete()
	default:
		return fmt.Errorf("route not specified for rpc %s at %s",
			method.Desc.FullName(), method.Location.SourceFile)
	}

	paramNames, pathPattern, ok := parsePath(info.Path.Raw)
	if !ok {
		return fmt.Errorf("invalid route path %s for rpc %s at %s",
			info.Path, method.Desc.FullName(), method.Location.SourceFile)
	}
	info.Path.Pattern = pathPattern
	if len(paramNames) != 0 {
		pbNameToField := make(map[string]*protogen.Field, len(method.Input.Fields))
		for _, v := range method.Input.Fields {
			pbNameToField[string(v.Desc.Name())] = v
		}
		paramRecord := make(map[string]struct{}, len(paramNames))
		info.PathParams = make([]*PathParamInfo, 0, len(paramNames))
		for _, paramName := range paramNames {
			if _, ok := paramRecord[paramName]; ok {
				return fmt.Errorf("repeated path parameter %s for rpc %s at %s",
					paramName, method.Desc.FullName(), method.Location.SourceFile)
			}
			field, ok := pbNameToField[paramName]
			if !ok {
				return fmt.Errorf("path param %s doesn't match any field name for rpc %s at %s",
					paramName, method.Desc.FullName(), method.Location.SourceFile)
			}
			paramInfo := &PathParamInfo{
				GoField:    field.GoName,
				PathParam:  paramName,
				RawComment: string(field.Comments.Leading),
			}
			switch field.Desc.Kind() {
			case protoreflect.BoolKind:
				paramInfo.ScalarType = ScalarTypeBool
			case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
				paramInfo.ScalarType = ScalarTypeInt32
			case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
				paramInfo.ScalarType = ScalarTypeUint32
			case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
				paramInfo.ScalarType = ScalarTypeInt64
			case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
				paramInfo.ScalarType = ScalarTypeUint64
			case protoreflect.FloatKind:
				paramInfo.ScalarType = ScalarTypeFloat32
			case protoreflect.DoubleKind:
				paramInfo.ScalarType = ScalarTypeFloat64
			case protoreflect.EnumKind:
				paramInfo.ScalarType = ScalarTypeEnum
				paramInfo.Enum = field.Enum
			case protoreflect.StringKind:
				paramInfo.ScalarType = ScalarTypeString
			default:
				return fmt.Errorf("the field which path parameter %s matches is not of basic type "+
					"for rpc %s at %s", paramName, method.Desc.FullName(), method.Location.SourceFile)
			}
			info.PathParams = append(info.PathParams, paramInfo)
			paramRecord[paramName] = struct{}{}
		}
	}
	info.AllInputFieldsInPath = len(paramNames) == len(method.Input.Fields)

	reg.rpcInfos[method] = &info
	reg.rpcList = append(reg.rpcList, method)
	return nil
}

func parsePath(path string) (_ []string, pattern string, _ bool) { // TODO: add unit test
	if path == "" {
		return nil, "", false
	}
	segments := strings.Split(path, "/")
	if segments[0] != "" {
		return nil, "", false
	}
	var params []string
	for k, v := range segments[1:] {
		if v == "" {
			if k == len(segments)-2 { // tailing slash
				return params, pattern + "/", true
			}
			return nil, "", false
		}
		if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
			param := v[1 : len(v)-1]
			if param == "" {
				return nil, "", false
			}
			if strings.ContainsAny(param, "{}") {
				return nil, "", false
			}
			params = append(params, param)
			pattern += "/{}"
			continue
		}
		if strings.ContainsAny(v, "{}") {
			return nil, "", false
		}
		pattern += "/" + v
	}
	return params, pattern, true
}
