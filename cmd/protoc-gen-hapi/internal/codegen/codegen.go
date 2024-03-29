package codegen

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/paleviews/hapi/cmd/protoc-gen-hapi/internal/serviceregistry"
	"github.com/paleviews/hapi/descriptor/annotations"
)

func Generate(reg *serviceregistry.Registry, plugin *protogen.Plugin, version string, noProtocVer bool) {
	for _, f := range plugin.Files {
		if f.Generate {
			codegen(reg, plugin, f, version, noProtocVer)
		}
	}
}

func codegen(
	reg *serviceregistry.Registry, plugin *protogen.Plugin, file *protogen.File,
	version string, noProtocVer bool,
) *protogen.GeneratedFile {

	cg := &codeGenerator{
		reg:         reg,
		plugin:      plugin,
		file:        file,
		g:           plugin.NewGeneratedFile(file.GeneratedFilenamePrefix+".hapi.go", file.GoImportPath),
		noProtocVer: noProtocVer,
	}

	cg.generateGoFileHeader(version)

	cg.g.P("package ", file.GoPackageName)
	cg.g.P()

	for _, v := range file.Enums {
		cg.generateEnum(v)
	}

	for _, v := range file.Messages {
		cg.generateMessage(v)
	}

	for _, v := range file.Services {
		cg.generateService(v)
	}

	return cg.g
}

type codeGenerator struct {
	reg         *serviceregistry.Registry
	plugin      *protogen.Plugin
	file        *protogen.File
	g           *protogen.GeneratedFile
	noProtocVer bool
}

func (cg *codeGenerator) generateGoFileHeader(version string) {
	cg.g.P("// Code generated by proto-gen-hapi. DO NOT EDIT.")
	cg.g.P("// version:")
	protocVersion := "(unknown)"
	if v := cg.plugin.Request.GetCompilerVersion(); v != nil {
		protocVersion = fmt.Sprintf("v%v.%v.%v", v.GetMajor(), v.GetMinor(), v.GetPatch())
		if s := v.GetSuffix(); s != "" {
			protocVersion += "-" + s
		}
	}
	cg.g.P("// - protoc-gen-hapi ", version)
	if !cg.noProtocVer {
		cg.g.P("// - protoc          ", protocVersion)
	}
	if cg.file.Proto.GetOptions().GetDeprecated() {
		cg.g.P("// ", cg.file.Desc.Path(), " is a deprecated file.")
	} else {
		cg.g.P("// source: ", cg.file.Desc.Path())
	}
	cg.g.P()
}

var (
	goIdentResponseCodeOKDefault              = hapiRuntimePackage.Ident("DefaultResponseCodeOK")
	goIdentResponseCodeInvalidInputDefault    = hapiRuntimePackage.Ident("DefaultResponseCodeInvalidInput")
	goIdentResponseCodeUnauthenticatedDefault = hapiRuntimePackage.Ident("DefaultResponseCodeUnauthenticated")
	goIdentResponseCodeServerErrorDefault     = hapiRuntimePackage.Ident("DefaultResponseCodeServerError")
	goIdentResponseCodeDescGetterDefault      = hapiRuntimePackage.Ident("GetDescFromResponseCode")
	goIdentAPIErrorFromResponseCodeDefault    = hapiRuntimePackage.Ident("APIErrorFromResponseCode")
)

func (cg *codeGenerator) getCodeOKIdent() protogen.GoIdent {
	info := cg.reg.ResponseCodeInfo()
	if info == nil {
		return goIdentResponseCodeOKDefault
	}
	return info.ValueOK.GoIdent
}

func (cg *codeGenerator) getCodeInvalidInputIdent() protogen.GoIdent {
	info := cg.reg.ResponseCodeInfo()
	if info == nil {
		return goIdentResponseCodeInvalidInputDefault
	}
	return info.ValueInvalidInput.GoIdent
}

func (cg *codeGenerator) getCodeUnauthenticatedIdent() protogen.GoIdent {
	info := cg.reg.ResponseCodeInfo()
	if info == nil {
		return goIdentResponseCodeUnauthenticatedDefault
	}
	return info.ValueUnauthenticated.GoIdent
}

func (cg *codeGenerator) getCodeServerErrorIdent() protogen.GoIdent {
	info := cg.reg.ResponseCodeInfo()
	if info == nil {
		return goIdentResponseCodeServerErrorDefault
	}
	return info.ValueServerError.GoIdent
}

func (cg *codeGenerator) getResponseCodeDescGetterIdent() protogen.GoIdent {
	info := cg.reg.ResponseCodeInfo()
	if info == nil {
		return goIdentResponseCodeDescGetterDefault
	}
	i := info.OriginEnum.GoIdent
	return i.GoImportPath.Ident("GetDescBy" + i.GoName)
}

func (cg *codeGenerator) getAPIErrorFromResponseCodeIdent() protogen.GoIdent {
	info := cg.reg.ResponseCodeInfo()
	if info == nil {
		return goIdentAPIErrorFromResponseCodeDefault
	}
	i := info.OriginEnum.GoIdent
	return i.GoImportPath.Ident("APIErrorFrom" + i.GoName)

}

func (cg *codeGenerator) generateEnum(e *protogen.Enum) {
	var isResponseCode bool
	if info := cg.reg.ResponseCodeInfo(); info != nil && info.OriginEnum == e {
		isResponseCode = true
	}
	leadingComments := cg.appendDeprecationSuffix(e.Comments.Leading,
		e.Desc.Options().(*descriptorpb.EnumOptions).GetDeprecated())
	if isResponseCode {
		cg.g.P(leadingComments, "type ", e.GoIdent, " = ", hapiRuntimePackage.Ident("ResponseCode"))
	} else {
		cg.g.P(leadingComments, "type ", e.GoIdent, " int32")
	}
	cg.g.P()
	cg.g.P("const (")
	for _, v := range e.Values {
		leadingComments := cg.appendDeprecationSuffix(v.Comments.Leading,
			v.Desc.Options().(*descriptorpb.EnumValueOptions).GetDeprecated())
		cg.g.P(leadingComments, v.GoIdent, " ", e.GoIdent, " = ", v.Desc.Number())
	}
	cg.g.P(")")
	cg.g.P()
	if !isResponseCode {
		return
	}

	cg.g.P("func ", cg.getResponseCodeDescGetterIdent(), "(code ", e.GoIdent, ") string {")
	cg.g.P("switch code {")
	for _, v := range e.Values {
		cg.g.P("case ", v.GoIdent, ":")
		cg.g.P(`return "`, cg.reg.ResponseCodeInfo().DescriptionByEnumValue[v], `"`)
	}
	cg.g.P("default:")
	cg.g.P(`return ""`)
	cg.g.P("}")
	cg.g.P("}")
	cg.g.P()

	cg.g.P("func ", cg.getAPIErrorFromResponseCodeIdent(), "(code ", e.GoIdent, ", src error) ",
		hapiRuntimePackage.Ident("APIError"), " {")
	cg.g.P("return ", hapiRuntimePackage.Ident("NewAPIError"), "(code, ",
		cg.getResponseCodeDescGetterIdent(), "(code), src)")
	cg.g.P("}")
	cg.g.P()
}

func (cg *codeGenerator) generateMessage(msg *protogen.Message) {
	if msg.Desc.IsMapEntry() {
		return
	}

	for _, v := range msg.Enums {
		cg.generateEnum(v)
	}

	for _, v := range msg.Messages {
		cg.generateMessage(v)
	}

	leadingComments := cg.appendDeprecationSuffix(msg.Comments.Leading,
		msg.Desc.Options().(*descriptorpb.MessageOptions).GetDeprecated())
	if len(msg.Fields) == 0 {
		cg.g.P(leadingComments, "type ", msg.GoIdent, " struct {}")
	} else {
		cg.g.P(leadingComments, "type ", msg.GoIdent, " struct {")
		for _, v := range msg.Fields {
			cg.generateMessageField(v)
		}
		cg.g.P("}")
	}

	if len(msg.Fields) != 0 {
		for _, v := range msg.Fields {
			cg.generateFieldGetter(v)
		}
	}

	cg.g.P()
}

func (cg *codeGenerator) generateMessageField(f *protogen.Field) {
	goType := cg.fieldGoType(f)
	leadingComments := cg.appendDeprecationSuffix(f.Comments.Leading,
		f.Desc.Options().(*descriptorpb.FieldOptions).GetDeprecated())
	cg.g.P(leadingComments, f.GoName, " ", goType,
		" `json:\"", f.Desc.Name(), ",omitempty\"`")
}

func (cg *codeGenerator) generateFieldGetter(f *protogen.Field) {
	goType := cg.fieldGoType(f)
	zero := "nil"
	if !f.Desc.IsMap() && !f.Desc.IsList() {
		switch f.Desc.Kind() {
		case protoreflect.BoolKind:
			zero = "false"
		case protoreflect.EnumKind:
			zero = "0" // TODO: enum zero lit
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
			protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
			protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
			protoreflect.FloatKind,
			protoreflect.DoubleKind:
			zero = "0"
		case protoreflect.StringKind:
			zero = `""`
		}
	}
	cg.g.P()
	cg.g.P("func (x *", f.Parent.GoIdent, ") Get", f.GoName, "() ", goType, " {") // TODO: f.Parent could be nil
	cg.g.P("if x != nil {")
	cg.g.P("return x.", f.GoName)
	cg.g.P("}")
	cg.g.P("return ", zero)
	cg.g.P("}")
}

func (cg *codeGenerator) fieldGoType(field *protogen.Field) (goType string) {
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		goType = "bool"
	case protoreflect.EnumKind:
		goType = cg.g.QualifiedGoIdent(field.Enum.GoIdent)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		goType = "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		goType = "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		goType = "int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		goType = "uint64"
	case protoreflect.FloatKind:
		goType = "float32"
	case protoreflect.DoubleKind:
		goType = "float64"
	case protoreflect.StringKind:
		goType = "string"
	case protoreflect.BytesKind:
		goType = "[]byte"
	case protoreflect.MessageKind, protoreflect.GroupKind:
		if field.Message.GoIdent.GoName == "Any" &&
			field.Message.GoIdent.GoImportPath == "github.com/paleviews/hapi/descriptor/types" {
			goType = cg.g.QualifiedGoIdent(field.Message.GoIdent)
		} else {
			goType = "*" + cg.g.QualifiedGoIdent(field.Message.GoIdent)
		}
	}
	switch {
	case field.Desc.IsList():
		return "[]" + goType
	case field.Desc.IsMap():
		keyType := cg.fieldGoType(field.Message.Fields[0])
		valType := cg.fieldGoType(field.Message.Fields[1])
		return fmt.Sprintf("map[%v]%v", keyType, valType)
	}
	return goType
}

func (cg *codeGenerator) generateService(svc *protogen.Service) {
	cg.generateServiceServer(svc)
	cg.generateServiceConstructor(svc)
}

func (cg *codeGenerator) generateServiceServer(svc *protogen.Service) {
	leadingComments := cg.appendDeprecationSuffix(svc.Comments.Leading,
		svc.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated())
	cg.g.P(leadingComments, "type ", svc.GoName, "Server interface {")
	for _, v := range svc.Methods {
		cg.generateSignature(v)
	}
	cg.g.P("}")
	cg.g.P()
}

const (
	contextPackage     = protogen.GoImportPath("context")
	hapiRuntimePackage = protogen.GoImportPath("github.com/paleviews/hapi/runtime")
	errorsPackage      = protogen.GoImportPath("errors")
	httpPackage        = protogen.GoImportPath("net/http")
	fmtPackage         = protogen.GoImportPath("fmt")
	ioPackage          = protogen.GoImportPath("io")
	base64Package      = protogen.GoImportPath("encoding/base64")
)

func (cg *codeGenerator) generateSignature(m *protogen.Method) {
	leadingComments := cg.appendDeprecationSuffix(m.Comments.Leading,
		m.Desc.Options().(*descriptorpb.MethodOptions).GetDeprecated())
	cg.g.P(leadingComments,
		m.GoName, "(", contextPackage.Ident("Context"), ", *", m.Input.GoIdent, ") (*",
		m.Output.GoIdent, ", error)")
}

const (
	serverVar              = "svr"
	handlerFacilitatorVar  = "hf"
	inputVar               = "input"
	hapiErrorVar           = "hapiError"
	responseWrapperVar     = "responseWrapper"
	responseWrapperFuncVar = "responseWrapperFunc"
	authMiddlewareVar      = "authMiddleware"
	formDecoderVar         = "decodeForm"
)

func (cg *codeGenerator) generateServiceConstructor(svc *protogen.Service) {
	cg.g.P("func New", svc.GoName, "Service(",
		serverVar, " ", svc.GoName, "Server, ", handlerFacilitatorVar, " ",
		hapiRuntimePackage.Ident("HandlerFacilitator"), ") ",
		hapiRuntimePackage.Ident("Service"), " {") // begin of constructor
	cg.g.P("if ", serverVar, " == nil {")
	cg.g.P("panic(", errorsPackage.Ident("New"), `("`, serverVar, ` is nil"))`)
	cg.g.P("}")
	cg.g.P("if ", handlerFacilitatorVar, " == nil {")
	cg.g.P("panic(", errorsPackage.Ident("New"), `("`, handlerFacilitatorVar, ` is nil"))`)
	cg.g.P("}")

	goIdentServerErrorCode := cg.getCodeServerErrorIdent()
	goIdentUnauthenticatedCode := cg.getCodeUnauthenticatedIdent()
	options := cg.reg.GlobalInfo()
	switch options.ResponseCodeIn {
	case annotations.ResponseCodeLocation_RESPONSE_CODE_LOCATION_BODY:
		cg.g.P(responseWrapperVar, " := ", hapiRuntimePackage.Ident("NewCodeInBodyWrapper"),
			"(", handlerFacilitatorVar, ", ", goIdentServerErrorCode, ",\n",
			cg.getResponseCodeDescGetterIdent(), "(", goIdentServerErrorCode, "))")
	case annotations.ResponseCodeLocation_RESPONSE_CODE_LOCATION_HEADER:
		cg.g.P(responseWrapperVar, " := ", hapiRuntimePackage.Ident("NewCodeInHeaderWrapper"),
			"(", handlerFacilitatorVar, ", ", goIdentServerErrorCode, ",\n",
			cg.getResponseCodeDescGetterIdent(), "(", goIdentServerErrorCode, "))")
	default:
		panic(serviceregistry.ErrUnreachableCode)
	}
	cg.g.P(responseWrapperFuncVar, " := ", responseWrapperVar, `.WriteResponse`)
	switch options.AuthKind {
	case annotations.AuthKind_AUTH_KIND_NONE:
	case annotations.AuthKind_AUTH_KIND_BEARER_IN_HEADER:
		cg.g.P(authMiddlewareVar, " := ", hapiRuntimePackage.Ident("AuthMiddleware"),
			"(", handlerFacilitatorVar, ", ", hapiRuntimePackage.Ident("GetBearerTokenInHeader"), ", \n",
			goIdentUnauthenticatedCode, ", \n",
			cg.getResponseCodeDescGetterIdent(), "(", goIdentUnauthenticatedCode, "), \n",
			responseWrapperVar, ")")
		cg.g.P("_ = ", authMiddlewareVar)
	default:
		panic(serviceregistry.ErrUnreachableCode)
	}
	cg.g.P()

	cg.g.P("return ", hapiRuntimePackage.Ident("Service"), "{") // begin of Service
	cg.g.P("RPCs: []", hapiRuntimePackage.Ident("RPC"), "{")    // begin of RPCs
	for _, v := range svc.Methods {
		cg.generateRPC(v)
	}
	cg.g.P("},") // end of RPCs
	cg.g.P("}")  // end of Service
	cg.g.P("}")  // end of constructor
}

func (cg *codeGenerator) generateRPC(m *protogen.Method) {
	info := cg.reg.RPCInfo(m)
	path := info.Path.Raw
	var methodIdent string
	switch info.HTTPMethod {
	case serviceregistry.HTTPMethodGet:
		methodIdent = "HTTPMethodGet"
	case serviceregistry.HTTPMethodPost:
		methodIdent = "HTTPMethodPost"
	case serviceregistry.HTTPMethodPut:
		methodIdent = "HTTPMethodPut"
	case serviceregistry.HTTPMethodPatch:
		methodIdent = "HTTPMethodPatch"
	case serviceregistry.HTTPMethodDelete:
		methodIdent = "HTTPMethodDelete"
	default:
		panic(serviceregistry.ErrUnreachableCode)
	}
	cg.g.P("{")
	cg.g.P("Route: ", hapiRuntimePackage.Ident("Route"), "{")
	cg.g.P("Method: ", hapiRuntimePackage.Ident(methodIdent), ",")
	cg.g.P(`Path: "`, path, `",`)
	cg.g.P("},")
	cg.generateHandler(m, info)
	cg.g.P("},")
}

func (cg *codeGenerator) inputFromForm(rpcInfo *serviceregistry.RPCInfo) bool {
	return rpcInfo.HTTPMethod == serviceregistry.HTTPMethodGet ||
		rpcInfo.HTTPMethod == serviceregistry.HTTPMethodDelete
}

func (cg *codeGenerator) generateHandler(method *protogen.Method, rpcInfo *serviceregistry.RPCInfo) {
	hasAuthMiddleware := cg.reg.GlobalInfo().AuthKind != annotations.AuthKind_AUTH_KIND_NONE && !rpcInfo.SkipAuth
	if hasAuthMiddleware {
		cg.g.P("Handler: ", authMiddlewareVar, "(func(rw ", httpPackage.Ident("ResponseWriter"),
			", req *", httpPackage.Ident("Request"), ") {")
	} else {
		cg.g.P("Handler: func(rw ", httpPackage.Ident("ResponseWriter"),
			", req *", httpPackage.Ident("Request"), ") {")
	}

	if !rpcInfo.AllInputFieldsInPath && cg.inputFromForm(rpcInfo) {
		cg.generateFormDecoder(method.Input, rpcInfo)
	}

	goIdentOKCode := cg.getCodeOKIdent()
	goIdentServerErrorCode := cg.getCodeServerErrorIdent()

	cg.g.P("ctx := req.Context()")
	cg.g.P("var (")
	cg.g.P(hapiErrorVar, " *", hapiRuntimePackage.Ident("APIError"))
	cg.g.P("data interface{}")
	cg.g.P(")")
	cg.g.P("defer func() {")
	cg.g.P("if ", hapiErrorVar, " == nil {")
	cg.g.P(handlerFacilitatorVar, ".ResultHook(ctx, data)")
	cg.g.P(responseWrapperFuncVar, "(ctx, rw, ", goIdentOKCode, `, "OK", data)`)
	cg.g.P("} else {")
	cg.g.P(handlerFacilitatorVar, ".ErrorHook(ctx, *", hapiErrorVar, ")")
	cg.g.P(responseWrapperFuncVar, "(ctx, rw, ", hapiErrorVar, ".Code, ", hapiErrorVar, ".Message, nil)")
	cg.g.P("}")
	cg.g.P("}()")

	cg.generateInputAssembleStatements(method, rpcInfo)

	cg.g.P("// invoke server logic")
	cg.g.P("output, err := ", serverVar, ".", method.GoName, "(ctx, ", inputVar, ")")
	cg.g.P("if err != nil {")
	cg.g.P("var hapiError1 ", hapiRuntimePackage.Ident("APIError"))
	cg.g.P("if ok := ", errorsPackage.Ident("As"), "(err, &hapiError1); !ok {")
	cg.g.P("hapiError1 = ", cg.getAPIErrorFromResponseCodeIdent(), "(")
	cg.g.P(goIdentServerErrorCode, ",")
	cg.g.P(fmtPackage.Ident("Errorf"), `("logic error: %w", err),`)
	cg.g.P(")")
	cg.g.P("}")
	cg.g.P(hapiErrorVar, " = &hapiError1")
	cg.g.P("return")
	cg.g.P("}")
	cg.g.P("data = output")
	if hasAuthMiddleware {
		cg.g.P("}),")
	} else {
		cg.g.P("},")
	}
}

func (cg *codeGenerator) generateFormDecoder(input *protogen.Message, rpcInfo *serviceregistry.RPCInfo) {
	skip := make(map[string]struct{}, len(rpcInfo.PathParams))
	for _, v := range rpcInfo.PathParams {
		skip[v.PathParam] = struct{}{}
	}

	cg.g.P(formDecoderVar, " := func(req *http.Request) (*",
		cg.g.QualifiedGoIdent(input.GoIdent), ", error) {") // +1
	cg.g.P("err := req.ParseForm()")
	cg.g.P("if err != nil {")
	cg.g.P("return nil, ", fmtPackage.Ident("Errorf"), `("parse form: %w", err)`)
	cg.g.P("}")
	cg.g.P("var input ", cg.g.QualifiedGoIdent(input.GoIdent))
	for _, v := range input.Fields {
		key := string(v.Desc.Name())
		if _, ok := skip[key]; ok {
			continue
		}
		lit := cg.fieldGoType(v)
		scalarParsePrint := func(parser string) {
			cg.g.P(`if s := req.Form.Get("`, key, `"); s != "" {`)
			cg.g.P("input.", v.GoName, ", err = ", hapiRuntimePackage.Ident(parser), "(s)")
			cg.g.P("if err != nil {")
			cg.g.P("return nil, ", fmtPackage.Ident("Errorf"), `("parse query `, key, ` as `, lit, `: %w", err)`)
			cg.g.P("}")
			cg.g.P("}")
		}
		switch lit {
		case "bool":
			scalarParsePrint("ParseBool")
		case "int32":
			scalarParsePrint("ParseInt32")
		case "int64":
			scalarParsePrint("ParseInt64")
		case "uint32":
			scalarParsePrint("ParseUint32")
		case "uint64":
			scalarParsePrint("ParseUint64")
		case "float32":
			scalarParsePrint("ParseFloat32")
		case "float64":
			scalarParsePrint("ParseFloat64")
		case "string":
			cg.g.P("input.", v.GoName, ` = req.Form.Get("`, string(v.Desc.Name()), `")`)
		case "[]byte":
			cg.g.P(`if s := req.Form.Get("`, key, `"); s != "" {`)
			cg.g.P("input.", v.GoName, ", err = ", base64Package.Ident("StdEncoding.DecodeString"), "(s)")
			cg.g.P("if err != nil {")
			cg.g.P("return nil, ", fmtPackage.Ident("Errorf"), `("parse query `, key, ` as `, lit, `: %w", err)`)
			cg.g.P("}")
			cg.g.P("}")
		default:
			switch {
			case lit[0] == '*', lit[:3] == "map":
				cg.g.P(`if s := req.Form.Get("`, key, `"); s != "" {`)
				cg.g.P(`err = hf.DecodeJSON([]byte(s), &input.`, v.GoName, ")")
				cg.g.P("if err != nil {")
				cg.g.P("return nil, ", fmtPackage.Ident("Errorf"), `("json decode query `, key, `: %w", err)`)
				cg.g.P("}")
				cg.g.P("}")
			case lit[:2] == "[]":
				lit := lit[2:]
				scalarParsePrint := func(parser string) {
					cg.g.P(`for _, v := range req.Form["`, key, `"] {`)
					cg.g.P("tmp, err := ", hapiRuntimePackage.Ident(parser), "(v)")
					cg.g.P("if err != nil {")
					cg.g.P("return nil, ", fmtPackage.Ident("Errorf"), `("parse query `, key, ` as []`, lit, `: %w", err)`)
					cg.g.P("}")
					cg.g.P("input.", v.GoName, " = append(input.", v.GoName, ", tmp)")
					cg.g.P("}")
				}
				switch lit {
				case "bool":
					scalarParsePrint("ParseBool")
				case "int32":
					scalarParsePrint("ParseInt32")
				case "int64":
					scalarParsePrint("ParseInt64")
				case "uint32":
					scalarParsePrint("ParseUint32")
				case "uint64":
					scalarParsePrint("ParseUint64")
				case "float32":
					scalarParsePrint("ParseFloat32")
				case "float64":
					scalarParsePrint("ParseFloat64")
				case "string":
					cg.g.P("input.", v.GoName, ` = req.Form["`, key, `"]`)
				case "[]byte":
					cg.g.P(`for _, v := range req.Form["`, key, `"] {`)
					cg.g.P("tmp, err := ", base64Package.Ident("StdEncoding.DecodeString"), "(v)")
					cg.g.P("if err != nil {")
					cg.g.P("return nil, ", fmtPackage.Ident("Errorf"), `("parse query `, key, ` as []`, lit, `: %w", err)`)
					cg.g.P("}")
					cg.g.P("input.", v.GoName, " = append(input.", v.GoName, ", tmp)")
					cg.g.P("}")
				default:
					switch {
					case lit[0] == '*':
						cg.g.P(`for _, v := range req.Form["`, key, `"] {`)
						cg.g.P(`var tmp `, lit)
						cg.g.P(`err := hf.DecodeJSON([]byte(v), &tmp)`)
						cg.g.P("if err != nil {")
						cg.g.P("return nil, ", fmtPackage.Ident("Errorf"), `("json decode query `, key, `: %w", err)`)
						cg.g.P("}")
						cg.g.P("input.", v.GoName, " = append(input.", v.GoName, ", tmp)")
						cg.g.P("}")
					case lit[:3] == "map":
						panic(fmt.Errorf("list of maps ([]%s) can't be defined in protobuf", lit))
					case lit[:2] == "[]":
						panic(fmt.Errorf("list of lists ([]%s) can't be defined in protobuf", lit))
					default:
						cg.g.P(`for _, v := range req.Form["`, key, `"] {`)
						cg.g.P("tmp, err := ", hapiRuntimePackage.Ident("ParseInt32"), "(v)")
						cg.g.P("if err != nil {")
						cg.g.P("return nil, ", fmtPackage.Ident("Errorf"), `("parse query `, key, ` as []`, lit, `: %w", err)`)
						cg.g.P("}")
						cg.g.P("input.", v.GoName, " = append(input.", v.GoName, ", ", lit, "(tmp))")
						cg.g.P("}")
					}
				}
			default:
				cg.g.P(`if s := req.Form.Get("`, key, `"); s != "" {`)
				cg.g.P("tmp, err := ", hapiRuntimePackage.Ident("ParseInt32"), "(s)")
				cg.g.P("if err != nil {")
				cg.g.P("return nil, ", fmtPackage.Ident("Errorf"), `("parse query `, key, ` as `, lit, `: %w", err)`)
				cg.g.P("}")
				cg.g.P("input.", v.GoName, " = ", lit, "(tmp)")
				cg.g.P("}")
			}
		}
	}
	cg.g.P("return &input, nil")
	cg.g.P("}") // -1
}

func (cg *codeGenerator) generateInputAssembleStatements(method *protogen.Method, rpcInfo *serviceregistry.RPCInfo) {
	input := method.Input
	fieldAssignments := cg.generatePathParamStatements(rpcInfo)
	if rpcInfo.AllInputFieldsInPath {
		if len(input.Fields) == 0 {
			cg.g.P(inputVar, " := new(", input.GoIdent, ")")
		} else {
			cg.g.P(inputVar, " := &", input.GoIdent, "{")
			for _, v := range fieldAssignments {
				cg.g.P(v.field, ": ", v._var, ",")
			}
			cg.g.P("}")
		}
	} else {
		goIdentServerErrorCode := cg.getCodeServerErrorIdent()
		goIdentInvalidInputCode := cg.getCodeInvalidInputIdent()
		cg.g.P("// decode and assemble input")
		if cg.inputFromForm(rpcInfo) {
			cg.g.P(inputVar, ", err := ", formDecoderVar, "(req)")
			cg.g.P("if err != nil {")
			cg.g.P("hapiError1 := ", cg.getAPIErrorFromResponseCodeIdent(), "(")
			cg.g.P(goIdentInvalidInputCode, ",")
			cg.g.P(fmtPackage.Ident("Errorf"), `("decode form: %w", err),`)
			cg.g.P(")")
			cg.g.P(hapiErrorVar, " = &hapiError1")
			cg.g.P("return")
			cg.g.P("}")
		} else {
			cg.g.P("bodyBytes, err := ", ioPackage.Ident("ReadAll"), "(req.Body)")
			cg.g.P("if err != nil {")
			cg.g.P("hapiError1 := ", cg.getAPIErrorFromResponseCodeIdent(), "(")
			cg.g.P(goIdentServerErrorCode, ",")
			cg.g.P(fmtPackage.Ident("Errorf"), `("read request body: %w", err),`)
			cg.g.P(")")
			cg.g.P(hapiErrorVar, " = &hapiError1")
			cg.g.P("return")
			cg.g.P("}")
			cg.g.P("var ", inputVar, " *", input.GoIdent)
			cg.g.P("err = ", handlerFacilitatorVar, ".DecodeJSON(bodyBytes, &", inputVar, ")")
			cg.g.P("if err != nil {")
			cg.g.P("hapiError1 := ", cg.getAPIErrorFromResponseCodeIdent(), "(")
			cg.g.P(goIdentInvalidInputCode, ",")
			cg.g.P(fmtPackage.Ident("Errorf"), `("decode json: %w", err),`)
			cg.g.P(")")
			cg.g.P(hapiErrorVar, " = &hapiError1")
			cg.g.P("return")
			cg.g.P("}")
		}
		for _, v := range fieldAssignments {
			cg.g.P(inputVar, ".", v.field, " = ", v._var)
		}
	}
}

type fieldAssignment struct {
	field string
	_var  string
}

func (cg *codeGenerator) generatePathParamStatements(rpcInfo *serviceregistry.RPCInfo) []fieldAssignment {
	goIdentInvalidInputCode := cg.getCodeInvalidInputIdent()
	fieldAssignments := make([]fieldAssignment, 0, len(rpcInfo.PathParams))
	for _, param := range rpcInfo.PathParams {
		strName := "str" + param.GoField
		cg.g.P("// get path parameter ", param.PathParam)
		cg.g.P(strName, ` := `, handlerFacilitatorVar, `.GetPathParam(req, "`, param.PathParam, `")`)
		cg.g.P("if ", strName, ` == "" {`)
		cg.g.P("hapiError1 := ", cg.getAPIErrorFromResponseCodeIdent(), "(")
		cg.g.P(goIdentInvalidInputCode, ",")
		cg.g.P(errorsPackage.Ident("New"), `("invalid `, param.PathParam, ` in path"),`)
		cg.g.P(")")
		cg.g.P(hapiErrorVar, " = &hapiError1")
		cg.g.P("return")
		cg.g.P("}")

		var (
			parser     interface{}
			realName   = "real" + param.GoField
			parsedName = realName
		)
		switch param.ScalarType {
		case serviceregistry.ScalarTypeBool:
			parser = hapiRuntimePackage.Ident("ParseBool")
		case serviceregistry.ScalarTypeInt32:
			parser = hapiRuntimePackage.Ident("ParseInt32")
		case serviceregistry.ScalarTypeUint32:
			parser = hapiRuntimePackage.Ident("ParseUint32")
		case serviceregistry.ScalarTypeInt64:
			parser = hapiRuntimePackage.Ident("ParseInt64")
		case serviceregistry.ScalarTypeUint64:
			parser = hapiRuntimePackage.Ident("ParseUint64")
		case serviceregistry.ScalarTypeFloat32:
			parser = hapiRuntimePackage.Ident("ParseFloat32")
		case serviceregistry.ScalarTypeFloat64:
			parser = hapiRuntimePackage.Ident("ParseFloat64")
		case serviceregistry.ScalarTypeEnum:
			parser = hapiRuntimePackage.Ident("ParseInt32")
			parsedName = "tmp" + param.GoField
		case serviceregistry.ScalarTypeString:
			realName = strName
		default:
			panic(serviceregistry.ErrUnreachableCode)
		}
		if parser != nil {
			cg.g.P(parsedName, ", err := ", parser, "(", strName, ")")
			cg.g.P("if err != nil {")
			cg.g.P("hapiError1 := ", cg.getAPIErrorFromResponseCodeIdent(), "(")
			cg.g.P(goIdentInvalidInputCode, ",")
			cg.g.P(fmtPackage.Ident("Errorf"), `("invalid `, param.PathParam, ` in path: %w", err),`)
			cg.g.P(")")
			cg.g.P(hapiErrorVar, " = &hapiError1")
			cg.g.P("return")
			cg.g.P("}")
			if realName != parsedName {
				cg.g.P(realName, " := ", cg.g.QualifiedGoIdent(param.Enum.GoIdent), "(", parsedName, ")")
			}
		}
		fieldAssignments = append(fieldAssignments, fieldAssignment{
			field: param.GoField,
			_var:  realName,
		})
	}
	return fieldAssignments
}

func (cg *codeGenerator) appendDeprecationSuffix(prefix protogen.Comments, deprecated bool) protogen.Comments {
	if !deprecated {
		return prefix
	}
	if prefix != "" {
		prefix += "\n"
	}
	return prefix + " Deprecated: Do not use.\n"
}
