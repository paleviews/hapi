// Code generated by proto-gen-hapi. DO NOT EDIT.
// version:
// - protoc-gen-hapi v0.1.0
// - protoc          v5.28.3
// source: todo/service.proto

package todo

import (
	context "context"
	errors "errors"
	fmt "fmt"
	types "github.com/paleviews/hapi/descriptor/types"
	codes "github.com/paleviews/hapi/example/todo/apidesign/golang/codes"
	common "github.com/paleviews/hapi/example/todo/apidesign/golang/common"
	runtime "github.com/paleviews/hapi/runtime"
	io "io"
	http "net/http"
)

type CreateRequest struct {
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func (x *CreateRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *CreateRequest) GetDetail() string {
	if x != nil {
		return x.Detail
	}
	return ""
}

type CreateResponse struct {
	ID int64 `json:"ID,omitempty"`
}

func (x *CreateResponse) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

type GetRequest struct {
	ID int64 `json:"ID,omitempty"`
}

func (x *GetRequest) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

type ListRequest struct {
	TitleContains  string `json:"title_contains,omitempty"`
	DetailContains string `json:"detail_contains,omitempty"`
	// start from 0
	Page     int64 `json:"page,omitempty"`
	PageSize int64 `json:"page_size,omitempty"`
}

func (x *ListRequest) GetTitleContains() string {
	if x != nil {
		return x.TitleContains
	}
	return ""
}

func (x *ListRequest) GetDetailContains() string {
	if x != nil {
		return x.DetailContains
	}
	return ""
}

func (x *ListRequest) GetPage() int64 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListRequest) GetPageSize() int64 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type ListResponse struct {
	Total int64   `json:"total,omitempty"`
	List  []*Todo `json:"list,omitempty"`
}

func (x *ListResponse) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *ListResponse) GetList() []*Todo {
	if x != nil {
		return x.List
	}
	return nil
}

type DeleteRequest struct {
	ID         int64                `json:"ID,omitempty"`
	SoftDelete bool                 `json:"soft_delete,omitempty"`
	More       map[string]types.Any `json:"more,omitempty"`
}

func (x *DeleteRequest) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *DeleteRequest) GetSoftDelete() bool {
	if x != nil {
		return x.SoftDelete
	}
	return false
}

func (x *DeleteRequest) GetMore() map[string]types.Any {
	if x != nil {
		return x.More
	}
	return nil
}

type V1Server interface {
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Get(context.Context, *GetRequest) (*Todo, error)
	List(context.Context, *ListRequest) (*ListResponse, error)
	Update(context.Context, *Todo) (*common.Empty, error)
	Delete(context.Context, *DeleteRequest) (*common.Empty, error)
}

func NewV1Service(svr V1Server, hf runtime.HandlerFacilitator) runtime.Service {
	if svr == nil {
		panic(errors.New("svr is nil"))
	}
	if hf == nil {
		panic(errors.New("hf is nil"))
	}
	responseWrapper := runtime.NewCodeInBodyWrapper(hf, codes.ResponseCode_RESPONSE_CODE_SERVER_ERROR,
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
					Method: runtime.HTTPMethodPost,
					Path:   "/todo/v1",
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
					// decode and assemble input
					bodyBytes, err := io.ReadAll(req.Body)
					if err != nil {
						hapiError1 := codes.APIErrorFromResponseCode(
							codes.ResponseCode_RESPONSE_CODE_SERVER_ERROR,
							fmt.Errorf("read request body: %w", err),
						)
						hapiError = &hapiError1
						return
					}
					var input *CreateRequest
					err = hf.DecodeJSON(bodyBytes, &input)
					if err != nil {
						hapiError1 := codes.APIErrorFromResponseCode(
							codes.ResponseCode_RESPONSE_CODE_INVALID_INPUT,
							fmt.Errorf("decode json: %w", err),
						)
						hapiError = &hapiError1
						return
					}
					// invoke server logic
					output, err := svr.Create(ctx, input)
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
					Method: runtime.HTTPMethodGet,
					Path:   "/todo/v1/{ID}",
				},
				Handler: func(rw http.ResponseWriter, req *http.Request) {
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
					// get path parameter ID
					strID := hf.GetPathParam(req, "ID")
					if strID == "" {
						hapiError1 := codes.APIErrorFromResponseCode(
							codes.ResponseCode_RESPONSE_CODE_INVALID_INPUT,
							errors.New("invalid ID in path"),
						)
						hapiError = &hapiError1
						return
					}
					realID, err := runtime.ParseInt64(strID)
					if err != nil {
						hapiError1 := codes.APIErrorFromResponseCode(
							codes.ResponseCode_RESPONSE_CODE_INVALID_INPUT,
							fmt.Errorf("invalid ID in path: %w", err),
						)
						hapiError = &hapiError1
						return
					}
					input := &GetRequest{
						ID: realID,
					}
					// invoke server logic
					output, err := svr.Get(ctx, input)
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
				},
			},
			{
				Route: runtime.Route{
					Method: runtime.HTTPMethodGet,
					Path:   "/todo/v1",
				},
				Handler: authMiddleware(func(rw http.ResponseWriter, req *http.Request) {
					decodeForm := func(req *http.Request) (*ListRequest, error) {
						err := req.ParseForm()
						if err != nil {
							return nil, fmt.Errorf("parse form: %w", err)
						}
						var input ListRequest
						input.TitleContains = req.Form.Get("title_contains")
						input.DetailContains = req.Form.Get("detail_contains")
						if s := req.Form.Get("page"); s != "" {
							input.Page, err = runtime.ParseInt64(s)
							if err != nil {
								return nil, fmt.Errorf("parse query page as int64: %w", err)
							}
						}
						if s := req.Form.Get("page_size"); s != "" {
							input.PageSize, err = runtime.ParseInt64(s)
							if err != nil {
								return nil, fmt.Errorf("parse query page_size as int64: %w", err)
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
					// invoke server logic
					output, err := svr.List(ctx, input)
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
					Method: runtime.HTTPMethodPut,
					Path:   "/todo/v1/{ID}",
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
					// get path parameter ID
					strID := hf.GetPathParam(req, "ID")
					if strID == "" {
						hapiError1 := codes.APIErrorFromResponseCode(
							codes.ResponseCode_RESPONSE_CODE_INVALID_INPUT,
							errors.New("invalid ID in path"),
						)
						hapiError = &hapiError1
						return
					}
					realID, err := runtime.ParseInt64(strID)
					if err != nil {
						hapiError1 := codes.APIErrorFromResponseCode(
							codes.ResponseCode_RESPONSE_CODE_INVALID_INPUT,
							fmt.Errorf("invalid ID in path: %w", err),
						)
						hapiError = &hapiError1
						return
					}
					// decode and assemble input
					bodyBytes, err := io.ReadAll(req.Body)
					if err != nil {
						hapiError1 := codes.APIErrorFromResponseCode(
							codes.ResponseCode_RESPONSE_CODE_SERVER_ERROR,
							fmt.Errorf("read request body: %w", err),
						)
						hapiError = &hapiError1
						return
					}
					var input *Todo
					err = hf.DecodeJSON(bodyBytes, &input)
					if err != nil {
						hapiError1 := codes.APIErrorFromResponseCode(
							codes.ResponseCode_RESPONSE_CODE_INVALID_INPUT,
							fmt.Errorf("decode json: %w", err),
						)
						hapiError = &hapiError1
						return
					}
					input.ID = realID
					// invoke server logic
					output, err := svr.Update(ctx, input)
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
					Path:   "/todo/v1/{ID}",
				},
				Handler: authMiddleware(func(rw http.ResponseWriter, req *http.Request) {
					decodeForm := func(req *http.Request) (*DeleteRequest, error) {
						err := req.ParseForm()
						if err != nil {
							return nil, fmt.Errorf("parse form: %w", err)
						}
						var input DeleteRequest
						if s := req.Form.Get("soft_delete"); s != "" {
							input.SoftDelete, err = runtime.ParseBool(s)
							if err != nil {
								return nil, fmt.Errorf("parse query soft_delete as bool: %w", err)
							}
						}
						if s := req.Form.Get("more"); s != "" {
							err = hf.DecodeJSON([]byte(s), &input.More)
							if err != nil {
								return nil, fmt.Errorf("json decode query more: %w", err)
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
					// get path parameter ID
					strID := hf.GetPathParam(req, "ID")
					if strID == "" {
						hapiError1 := codes.APIErrorFromResponseCode(
							codes.ResponseCode_RESPONSE_CODE_INVALID_INPUT,
							errors.New("invalid ID in path"),
						)
						hapiError = &hapiError1
						return
					}
					realID, err := runtime.ParseInt64(strID)
					if err != nil {
						hapiError1 := codes.APIErrorFromResponseCode(
							codes.ResponseCode_RESPONSE_CODE_INVALID_INPUT,
							fmt.Errorf("invalid ID in path: %w", err),
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
					input.ID = realID
					// invoke server logic
					output, err := svr.Delete(ctx, input)
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
