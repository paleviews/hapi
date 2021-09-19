package main

import (
	"net/http"

	"github.com/go-chi/cors"

	"github.com/paleviews/hapi/example/todo/apidesign/golang/todo"
	"github.com/paleviews/hapi/example/todo/handler"
	"github.com/paleviews/hapi/example/todo/logic"
	"github.com/paleviews/hapi/runtime"
)

func main() {
	hf := handler.NewFacilitator()
	service := todo.NewV1Service(logic.NewTodo(), hf)
	mux := hf.Mux()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut,
			http.MethodDelete, http.MethodOptions, http.MethodPatch},
		AllowedHeaders:   []string{"Origin", "Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link", "Content-Type", "X-Request-Id"},
		AllowCredentials: false,
		MaxAge:           3600,
	}))
	for _, rpc := range service.RPCs {
		switch rpc.Route.Method {
		case runtime.HTTPMethodGet:
			mux.Get(rpc.Route.Path, rpc.Handler)
		case runtime.HTTPMethodPost:
			mux.Post(rpc.Route.Path, rpc.Handler)
		case runtime.HTTPMethodPut:
			mux.Put(rpc.Route.Path, rpc.Handler)
		case runtime.HTTPMethodPatch:
			mux.Patch(rpc.Route.Path, rpc.Handler)
		case runtime.HTTPMethodDelete:
			mux.Delete(rpc.Route.Path, rpc.Handler)
		}
	}
	_ = http.ListenAndServe(":8080", mux)
}
