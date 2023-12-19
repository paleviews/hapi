// Code generated by proto-gen-hapi. DO NOT EDIT.
// version:
// - protoc-gen-hapi v0.1.0
// - protoc          v4.25.1
// source: testcase/service.proto

package testcase

import (
	context "context"
	base64 "encoding/base64"
	errors "errors"
	fmt "fmt"
	codes "github.com/paleviews/hapi/example/testcase/apidesign/golang/codes"
	common "github.com/paleviews/hapi/example/testcase/apidesign/golang/common"
	runtime "github.com/paleviews/hapi/runtime"
	http "net/http"
)

// Obj nested message comment
type DeepQueryRequest_Obj struct {
	// parent nested field comment
	Parent *DeepQueryRequest `json:"parent,omitempty"`
	Ratio  float64           `json:"ratio,omitempty"`
}

func (x *DeepQueryRequest_Obj) GetParent() *DeepQueryRequest {
	if x != nil {
		return x.Parent
	}
	return nil
}

func (x *DeepQueryRequest_Obj) GetRatio() float64 {
	if x != nil {
		return x.Ratio
	}
	return 0
}

type DeepQueryRequest_RefOnce struct {
	// id is id
	Id string `json:"id,omitempty"`
	// num is num
	Num int32 `json:"num,omitempty"`
}

func (x *DeepQueryRequest_RefOnce) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *DeepQueryRequest_RefOnce) GetNum() int32 {
	if x != nil {
		return x.Num
	}
	return 0
}

// DeepQueryRequest message comment
type DeepQueryRequest struct {
	// uuid field comment
	Uuid          string           `json:"uuid,omitempty"`
	StringToInt64 map[string]int64 `json:"string_to_int64,omitempty"`
	// string_to_obj field comment on map
	StringToObj map[string]*DeepQueryRequest_Obj `json:"string_to_obj,omitempty"`
	ObjsArray   []*DeepQueryRequest_Obj          `json:"objs_array,omitempty"`
	Obj         *DeepQueryRequest_Obj            `json:"obj,omitempty"`
	IsMarked    bool                             `json:"is_marked,omitempty"`
	// ref_once is ref_once
	RefOnce *DeepQueryRequest_RefOnce `json:"ref_once,omitempty"`
}

func (x *DeepQueryRequest) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *DeepQueryRequest) GetStringToInt64() map[string]int64 {
	if x != nil {
		return x.StringToInt64
	}
	return nil
}

func (x *DeepQueryRequest) GetStringToObj() map[string]*DeepQueryRequest_Obj {
	if x != nil {
		return x.StringToObj
	}
	return nil
}

func (x *DeepQueryRequest) GetObjsArray() []*DeepQueryRequest_Obj {
	if x != nil {
		return x.ObjsArray
	}
	return nil
}

func (x *DeepQueryRequest) GetObj() *DeepQueryRequest_Obj {
	if x != nil {
		return x.Obj
	}
	return nil
}

func (x *DeepQueryRequest) GetIsMarked() bool {
	if x != nil {
		return x.IsMarked
	}
	return false
}

func (x *DeepQueryRequest) GetRefOnce() *DeepQueryRequest_RefOnce {
	if x != nil {
		return x.RefOnce
	}
	return nil
}

// enum out
type DeepQueryResponse_Direction int32

const (
	// one line enum
	DeepQueryResponse_DIRECTION_UNKNOWN DeepQueryResponse_Direction = 0
	DeepQueryResponse_EAST              DeepQueryResponse_Direction = 1
	DeepQueryResponse_WEST              DeepQueryResponse_Direction = 2
	// multiple lines 1
	// multiple lines 2
	// multiple lines 3
	DeepQueryResponse_SOUTH DeepQueryResponse_Direction = 3
	DeepQueryResponse_NORTH DeepQueryResponse_Direction = 4
)

type DeepQueryResponse_RefOnce struct {
	// id is id
	Id string `json:"id,omitempty"`
	// num is num
	Num int32 `json:"num,omitempty"`
}

func (x *DeepQueryResponse_RefOnce) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *DeepQueryResponse_RefOnce) GetNum() int32 {
	if x != nil {
		return x.Num
	}
	return 0
}

type DeepQueryResponse struct {
	// ref_once in response
	RefOnce *DeepQueryResponse_RefOnce `json:"ref_once,omitempty"`
	// bs is bs
	Bs []byte `json:"bs,omitempty"`
	// direction line 1
	// line 2
	// line 3
	Direction DeepQueryResponse_Direction `json:"direction,omitempty"`
}

func (x *DeepQueryResponse) GetRefOnce() *DeepQueryResponse_RefOnce {
	if x != nil {
		return x.RefOnce
	}
	return nil
}

func (x *DeepQueryResponse) GetBs() []byte {
	if x != nil {
		return x.Bs
	}
	return nil
}

func (x *DeepQueryResponse) GetDirection() DeepQueryResponse_Direction {
	if x != nil {
		return x.Direction
	}
	return 0
}

type FormDecodeRequest_UUID int32

const (
	FormDecodeRequest_UUID_SHORT  FormDecodeRequest_UUID = 0
	FormDecodeRequest_UUID_MEDIUM FormDecodeRequest_UUID = 1
	FormDecodeRequest_UUID_LONG   FormDecodeRequest_UUID = 2
)

type FormDecodeRequest_Nested struct {
	Hello string `json:"hello,omitempty"`
	World bool   `json:"world,omitempty"`
}

func (x *FormDecodeRequest_Nested) GetHello() string {
	if x != nil {
		return x.Hello
	}
	return ""
}

func (x *FormDecodeRequest_Nested) GetWorld() bool {
	if x != nil {
		return x.World
	}
	return false
}

type FormDecodeRequest struct {
	Uuid         FormDecodeRequest_UUID               `json:"uuid,omitempty"`
	BoolField    bool                                 `json:"bool_field,omitempty"`
	Int32Field   int32                                `json:"int32_field,omitempty"`
	Int64Field   int64                                `json:"int64_field,omitempty"`
	Uint32Field  uint32                               `json:"uint32_field,omitempty"`
	Uint64Field  uint64                               `json:"uint64_field,omitempty"`
	Float32Field float32                              `json:"float32_field,omitempty"`
	Float64Field float64                              `json:"float64_field,omitempty"`
	StringField  string                               `json:"string_field,omitempty"`
	BytesField   []byte                               `json:"bytes_field,omitempty"`
	EnumField    FormDecodeRequest_UUID               `json:"enum_field,omitempty"`
	MessageField *FormDecodeRequest_Nested            `json:"message_field,omitempty"`
	SimpleMap    map[string]uint64                    `json:"simple_map,omitempty"`
	NotSimpleMap map[string]*FormDecodeRequest_Nested `json:"not_simple_map,omitempty"`
	BoolArray    []bool                               `json:"bool_array,omitempty"`
	Int32Array   []int32                              `json:"int32_array,omitempty"`
	Int64Array   []int64                              `json:"int64_array,omitempty"`
	Uint32Array  []uint32                             `json:"uint32_array,omitempty"`
	Uint64Array  []uint64                             `json:"uint64_array,omitempty"`
	Float32Array []float32                            `json:"float32_array,omitempty"`
	Float64Array []float64                            `json:"float64_array,omitempty"`
	StringArray  []string                             `json:"string_array,omitempty"`
	BytesArray   [][]byte                             `json:"bytes_array,omitempty"`
	EnumArray    []FormDecodeRequest_UUID             `json:"enum_array,omitempty"`
	MessageArray []*FormDecodeRequest_Nested          `json:"message_array,omitempty"`
}

func (x *FormDecodeRequest) GetUuid() FormDecodeRequest_UUID {
	if x != nil {
		return x.Uuid
	}
	return 0
}

func (x *FormDecodeRequest) GetBoolField() bool {
	if x != nil {
		return x.BoolField
	}
	return false
}

func (x *FormDecodeRequest) GetInt32Field() int32 {
	if x != nil {
		return x.Int32Field
	}
	return 0
}

func (x *FormDecodeRequest) GetInt64Field() int64 {
	if x != nil {
		return x.Int64Field
	}
	return 0
}

func (x *FormDecodeRequest) GetUint32Field() uint32 {
	if x != nil {
		return x.Uint32Field
	}
	return 0
}

func (x *FormDecodeRequest) GetUint64Field() uint64 {
	if x != nil {
		return x.Uint64Field
	}
	return 0
}

func (x *FormDecodeRequest) GetFloat32Field() float32 {
	if x != nil {
		return x.Float32Field
	}
	return 0
}

func (x *FormDecodeRequest) GetFloat64Field() float64 {
	if x != nil {
		return x.Float64Field
	}
	return 0
}

func (x *FormDecodeRequest) GetStringField() string {
	if x != nil {
		return x.StringField
	}
	return ""
}

func (x *FormDecodeRequest) GetBytesField() []byte {
	if x != nil {
		return x.BytesField
	}
	return nil
}

func (x *FormDecodeRequest) GetEnumField() FormDecodeRequest_UUID {
	if x != nil {
		return x.EnumField
	}
	return 0
}

func (x *FormDecodeRequest) GetMessageField() *FormDecodeRequest_Nested {
	if x != nil {
		return x.MessageField
	}
	return nil
}

func (x *FormDecodeRequest) GetSimpleMap() map[string]uint64 {
	if x != nil {
		return x.SimpleMap
	}
	return nil
}

func (x *FormDecodeRequest) GetNotSimpleMap() map[string]*FormDecodeRequest_Nested {
	if x != nil {
		return x.NotSimpleMap
	}
	return nil
}

func (x *FormDecodeRequest) GetBoolArray() []bool {
	if x != nil {
		return x.BoolArray
	}
	return nil
}

func (x *FormDecodeRequest) GetInt32Array() []int32 {
	if x != nil {
		return x.Int32Array
	}
	return nil
}

func (x *FormDecodeRequest) GetInt64Array() []int64 {
	if x != nil {
		return x.Int64Array
	}
	return nil
}

func (x *FormDecodeRequest) GetUint32Array() []uint32 {
	if x != nil {
		return x.Uint32Array
	}
	return nil
}

func (x *FormDecodeRequest) GetUint64Array() []uint64 {
	if x != nil {
		return x.Uint64Array
	}
	return nil
}

func (x *FormDecodeRequest) GetFloat32Array() []float32 {
	if x != nil {
		return x.Float32Array
	}
	return nil
}

func (x *FormDecodeRequest) GetFloat64Array() []float64 {
	if x != nil {
		return x.Float64Array
	}
	return nil
}

func (x *FormDecodeRequest) GetStringArray() []string {
	if x != nil {
		return x.StringArray
	}
	return nil
}

func (x *FormDecodeRequest) GetBytesArray() [][]byte {
	if x != nil {
		return x.BytesArray
	}
	return nil
}

func (x *FormDecodeRequest) GetEnumArray() []FormDecodeRequest_UUID {
	if x != nil {
		return x.EnumArray
	}
	return nil
}

func (x *FormDecodeRequest) GetMessageArray() []*FormDecodeRequest_Nested {
	if x != nil {
		return x.MessageArray
	}
	return nil
}

// V1 service comment
type V1Server interface {
	// DeepQuery rpc comment
	DeepQuery(context.Context, *DeepQueryRequest) (*DeepQueryResponse, error)
	FormDecode(context.Context, *FormDecodeRequest) (*common.Empty, error)
}

func NewV1Service(svr V1Server, hf runtime.HandlerFacilitator) runtime.Service {
	if svr == nil {
		panic(errors.New("svr is nil"))
	}
	if hf == nil {
		panic(errors.New("hf is nil"))
	}
	responseWrapper := runtime.NewCodeInHeaderWrapper(hf, codes.ResponseCode_RESPONSE_CODE_SERVER_ERROR,
		codes.GetDescByResponseCode(codes.ResponseCode_RESPONSE_CODE_SERVER_ERROR))
	responseWrapperFunc := responseWrapper.WriteResponse
	authMiddleware := runtime.AuthMiddleware(hf, runtime.GetBearerTokenInHeader,
		codes.ResponseCode_RESPONSE_CODE_UNAUTHENTICATED,
		codes.GetDescByResponseCode(codes.ResponseCode_RESPONSE_CODE_UNAUTHENTICATED),
		responseWrapper)
	_ = authMiddleware

	return runtime.Service{
		RPCs: []runtime.RPC{
			{
				Route: runtime.Route{
					Method: runtime.HTTPMethodGet,
					Path:   "/testcase/v1/{uuid}",
				},
				Handler: authMiddleware(func(rw http.ResponseWriter, req *http.Request) {
					decodeForm := func(req *http.Request) (*DeepQueryRequest, error) {
						err := req.ParseForm()
						if err != nil {
							return nil, fmt.Errorf("parse form: %w", err)
						}
						var input DeepQueryRequest
						if s := req.Form.Get("string_to_int64"); s != "" {
							err = hf.DecodeJSON([]byte(s), &input.StringToInt64)
							if err != nil {
								return nil, fmt.Errorf("json decode query string_to_int64: %w", err)
							}
						}
						if s := req.Form.Get("string_to_obj"); s != "" {
							err = hf.DecodeJSON([]byte(s), &input.StringToObj)
							if err != nil {
								return nil, fmt.Errorf("json decode query string_to_obj: %w", err)
							}
						}
						for _, v := range req.Form["objs_array"] {
							var tmp *DeepQueryRequest_Obj
							err := hf.DecodeJSON([]byte(v), &tmp)
							if err != nil {
								return nil, fmt.Errorf("json decode query objs_array: %w", err)
							}
							input.ObjsArray = append(input.ObjsArray, tmp)
						}
						if s := req.Form.Get("obj"); s != "" {
							err = hf.DecodeJSON([]byte(s), &input.Obj)
							if err != nil {
								return nil, fmt.Errorf("json decode query obj: %w", err)
							}
						}
						if s := req.Form.Get("is_marked"); s != "" {
							input.IsMarked, err = runtime.ParseBool(s)
							if err != nil {
								return nil, fmt.Errorf("parse query is_marked as bool: %w", err)
							}
						}
						if s := req.Form.Get("ref_once"); s != "" {
							err = hf.DecodeJSON([]byte(s), &input.RefOnce)
							if err != nil {
								return nil, fmt.Errorf("json decode query ref_once: %w", err)
							}
						}
						return &input, nil
					}
					ctx := req.Context()
					var (
						hapiError *runtime.APIError
						data      interface{}
					)
					defer func() {
						if hapiError == nil {
							hf.ResultHook(ctx, data)
							responseWrapperFunc(ctx, rw, codes.ResponseCode_RESPONSE_CODE_OK, "OK", data)
						} else {
							hf.ErrorHook(ctx, *hapiError)
							responseWrapperFunc(ctx, rw, hapiError.Code, hapiError.Message, nil)
						}
					}()
					// get path parameter uuid
					strUuid := hf.GetPathParam(req, "uuid")
					if strUuid == "" {
						hapiError1 := codes.APIErrorFromResponseCode(
							codes.ResponseCode_RESPONSE_CODE_INVALID_INPUT,
							errors.New("invalid uuid in path"),
						)
						hapiError = &hapiError1
						return
					}
					// decode and assemble input
					input, err := decodeForm(req)
					if err != nil {
						hapiError1 := codes.APIErrorFromResponseCode(
							codes.ResponseCode_RESPONSE_CODE_INVALID_INPUT,
							fmt.Errorf("decode form: %w", err),
						)
						hapiError = &hapiError1
						return
					}
					input.Uuid = strUuid
					// invoke server logic
					output, err := svr.DeepQuery(ctx, input)
					if err != nil {
						var hapiError1 runtime.APIError
						if ok := errors.As(err, &hapiError1); !ok {
							hapiError1 = codes.APIErrorFromResponseCode(
								codes.ResponseCode_RESPONSE_CODE_SERVER_ERROR,
								fmt.Errorf("logic error: %w", err),
							)
						}
						hapiError = &hapiError1
						return
					}
					data = output
				}),
			},
			{
				Route: runtime.Route{
					Method: runtime.HTTPMethodDelete,
					Path:   "/testcase/v1/{uuid}",
				},
				Handler: authMiddleware(func(rw http.ResponseWriter, req *http.Request) {
					decodeForm := func(req *http.Request) (*FormDecodeRequest, error) {
						err := req.ParseForm()
						if err != nil {
							return nil, fmt.Errorf("parse form: %w", err)
						}
						var input FormDecodeRequest
						if s := req.Form.Get("bool_field"); s != "" {
							input.BoolField, err = runtime.ParseBool(s)
							if err != nil {
								return nil, fmt.Errorf("parse query bool_field as bool: %w", err)
							}
						}
						if s := req.Form.Get("int32_field"); s != "" {
							input.Int32Field, err = runtime.ParseInt32(s)
							if err != nil {
								return nil, fmt.Errorf("parse query int32_field as int32: %w", err)
							}
						}
						if s := req.Form.Get("int64_field"); s != "" {
							input.Int64Field, err = runtime.ParseInt64(s)
							if err != nil {
								return nil, fmt.Errorf("parse query int64_field as int64: %w", err)
							}
						}
						if s := req.Form.Get("uint32_field"); s != "" {
							input.Uint32Field, err = runtime.ParseUint32(s)
							if err != nil {
								return nil, fmt.Errorf("parse query uint32_field as uint32: %w", err)
							}
						}
						if s := req.Form.Get("uint64_field"); s != "" {
							input.Uint64Field, err = runtime.ParseUint64(s)
							if err != nil {
								return nil, fmt.Errorf("parse query uint64_field as uint64: %w", err)
							}
						}
						if s := req.Form.Get("float32_field"); s != "" {
							input.Float32Field, err = runtime.ParseFloat32(s)
							if err != nil {
								return nil, fmt.Errorf("parse query float32_field as float32: %w", err)
							}
						}
						if s := req.Form.Get("float64_field"); s != "" {
							input.Float64Field, err = runtime.ParseFloat64(s)
							if err != nil {
								return nil, fmt.Errorf("parse query float64_field as float64: %w", err)
							}
						}
						input.StringField = req.Form.Get("string_field")
						if s := req.Form.Get("bytes_field"); s != "" {
							input.BytesField, err = base64.StdEncoding.DecodeString(s)
							if err != nil {
								return nil, fmt.Errorf("parse query bytes_field as []byte: %w", err)
							}
						}
						if s := req.Form.Get("enum_field"); s != "" {
							tmp, err := runtime.ParseInt32(s)
							if err != nil {
								return nil, fmt.Errorf("parse query enum_field as FormDecodeRequest_UUID: %w", err)
							}
							input.EnumField = FormDecodeRequest_UUID(tmp)
						}
						if s := req.Form.Get("message_field"); s != "" {
							err = hf.DecodeJSON([]byte(s), &input.MessageField)
							if err != nil {
								return nil, fmt.Errorf("json decode query message_field: %w", err)
							}
						}
						if s := req.Form.Get("simple_map"); s != "" {
							err = hf.DecodeJSON([]byte(s), &input.SimpleMap)
							if err != nil {
								return nil, fmt.Errorf("json decode query simple_map: %w", err)
							}
						}
						if s := req.Form.Get("not_simple_map"); s != "" {
							err = hf.DecodeJSON([]byte(s), &input.NotSimpleMap)
							if err != nil {
								return nil, fmt.Errorf("json decode query not_simple_map: %w", err)
							}
						}
						for _, v := range req.Form["bool_array"] {
							tmp, err := runtime.ParseBool(v)
							if err != nil {
								return nil, fmt.Errorf("parse query bool_array as []bool: %w", err)
							}
							input.BoolArray = append(input.BoolArray, tmp)
						}
						for _, v := range req.Form["int32_array"] {
							tmp, err := runtime.ParseInt32(v)
							if err != nil {
								return nil, fmt.Errorf("parse query int32_array as []int32: %w", err)
							}
							input.Int32Array = append(input.Int32Array, tmp)
						}
						for _, v := range req.Form["int64_array"] {
							tmp, err := runtime.ParseInt64(v)
							if err != nil {
								return nil, fmt.Errorf("parse query int64_array as []int64: %w", err)
							}
							input.Int64Array = append(input.Int64Array, tmp)
						}
						for _, v := range req.Form["uint32_array"] {
							tmp, err := runtime.ParseUint32(v)
							if err != nil {
								return nil, fmt.Errorf("parse query uint32_array as []uint32: %w", err)
							}
							input.Uint32Array = append(input.Uint32Array, tmp)
						}
						for _, v := range req.Form["uint64_array"] {
							tmp, err := runtime.ParseUint64(v)
							if err != nil {
								return nil, fmt.Errorf("parse query uint64_array as []uint64: %w", err)
							}
							input.Uint64Array = append(input.Uint64Array, tmp)
						}
						for _, v := range req.Form["float32_array"] {
							tmp, err := runtime.ParseFloat32(v)
							if err != nil {
								return nil, fmt.Errorf("parse query float32_array as []float32: %w", err)
							}
							input.Float32Array = append(input.Float32Array, tmp)
						}
						for _, v := range req.Form["float64_array"] {
							tmp, err := runtime.ParseFloat64(v)
							if err != nil {
								return nil, fmt.Errorf("parse query float64_array as []float64: %w", err)
							}
							input.Float64Array = append(input.Float64Array, tmp)
						}
						input.StringArray = req.Form["string_array"]
						for _, v := range req.Form["bytes_array"] {
							tmp, err := base64.StdEncoding.DecodeString(v)
							if err != nil {
								return nil, fmt.Errorf("parse query bytes_array as [][]byte: %w", err)
							}
							input.BytesArray = append(input.BytesArray, tmp)
						}
						for _, v := range req.Form["enum_array"] {
							tmp, err := runtime.ParseInt32(v)
							if err != nil {
								return nil, fmt.Errorf("parse query enum_array as []FormDecodeRequest_UUID: %w", err)
							}
							input.EnumArray = append(input.EnumArray, FormDecodeRequest_UUID(tmp))
						}
						for _, v := range req.Form["message_array"] {
							var tmp *FormDecodeRequest_Nested
							err := hf.DecodeJSON([]byte(v), &tmp)
							if err != nil {
								return nil, fmt.Errorf("json decode query message_array: %w", err)
							}
							input.MessageArray = append(input.MessageArray, tmp)
						}
						return &input, nil
					}
					ctx := req.Context()
					var (
						hapiError *runtime.APIError
						data      interface{}
					)
					defer func() {
						if hapiError == nil {
							hf.ResultHook(ctx, data)
							responseWrapperFunc(ctx, rw, codes.ResponseCode_RESPONSE_CODE_OK, "OK", data)
						} else {
							hf.ErrorHook(ctx, *hapiError)
							responseWrapperFunc(ctx, rw, hapiError.Code, hapiError.Message, nil)
						}
					}()
					// get path parameter uuid
					strUuid := hf.GetPathParam(req, "uuid")
					if strUuid == "" {
						hapiError1 := codes.APIErrorFromResponseCode(
							codes.ResponseCode_RESPONSE_CODE_INVALID_INPUT,
							errors.New("invalid uuid in path"),
						)
						hapiError = &hapiError1
						return
					}
					tmpUuid, err := runtime.ParseInt32(strUuid)
					if err != nil {
						hapiError1 := codes.APIErrorFromResponseCode(
							codes.ResponseCode_RESPONSE_CODE_INVALID_INPUT,
							fmt.Errorf("invalid uuid in path: %w", err),
						)
						hapiError = &hapiError1
						return
					}
					realUuid := FormDecodeRequest_UUID(tmpUuid)
					// decode and assemble input
					input, err := decodeForm(req)
					if err != nil {
						hapiError1 := codes.APIErrorFromResponseCode(
							codes.ResponseCode_RESPONSE_CODE_INVALID_INPUT,
							fmt.Errorf("decode form: %w", err),
						)
						hapiError = &hapiError1
						return
					}
					input.Uuid = realUuid
					// invoke server logic
					output, err := svr.FormDecode(ctx, input)
					if err != nil {
						var hapiError1 runtime.APIError
						if ok := errors.As(err, &hapiError1); !ok {
							hapiError1 = codes.APIErrorFromResponseCode(
								codes.ResponseCode_RESPONSE_CODE_SERVER_ERROR,
								fmt.Errorf("logic error: %w", err),
							)
						}
						hapiError = &hapiError1
						return
					}
					data = output
				}),
			},
		},
	}
}

type V2Server interface {
	CodeInHeaders(context.Context, *common.Empty) (*common.Empty, error)
}

func NewV2Service(svr V2Server, hf runtime.HandlerFacilitator) runtime.Service {
	if svr == nil {
		panic(errors.New("svr is nil"))
	}
	if hf == nil {
		panic(errors.New("hf is nil"))
	}
	responseWrapper := runtime.NewCodeInHeaderWrapper(hf, codes.ResponseCode_RESPONSE_CODE_SERVER_ERROR,
		codes.GetDescByResponseCode(codes.ResponseCode_RESPONSE_CODE_SERVER_ERROR))
	responseWrapperFunc := responseWrapper.WriteResponse
	authMiddleware := runtime.AuthMiddleware(hf, runtime.GetBearerTokenInHeader,
		codes.ResponseCode_RESPONSE_CODE_UNAUTHENTICATED,
		codes.GetDescByResponseCode(codes.ResponseCode_RESPONSE_CODE_UNAUTHENTICATED),
		responseWrapper)
	_ = authMiddleware

	return runtime.Service{
		RPCs: []runtime.RPC{
			{
				Route: runtime.Route{
					Method: runtime.HTTPMethodPut,
					Path:   "/testcase/v2",
				},
				Handler: authMiddleware(func(rw http.ResponseWriter, req *http.Request) {
					ctx := req.Context()
					var (
						hapiError *runtime.APIError
						data      interface{}
					)
					defer func() {
						if hapiError == nil {
							hf.ResultHook(ctx, data)
							responseWrapperFunc(ctx, rw, codes.ResponseCode_RESPONSE_CODE_OK, "OK", data)
						} else {
							hf.ErrorHook(ctx, *hapiError)
							responseWrapperFunc(ctx, rw, hapiError.Code, hapiError.Message, nil)
						}
					}()
					input := new(common.Empty)
					// invoke server logic
					output, err := svr.CodeInHeaders(ctx, input)
					if err != nil {
						var hapiError1 runtime.APIError
						if ok := errors.As(err, &hapiError1); !ok {
							hapiError1 = codes.APIErrorFromResponseCode(
								codes.ResponseCode_RESPONSE_CODE_SERVER_ERROR,
								fmt.Errorf("logic error: %w", err),
							)
						}
						hapiError = &hapiError1
						return
					}
					data = output
				}),
			},
		},
	}
}
