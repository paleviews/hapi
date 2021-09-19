module github.com/paleviews/hapi/example/todo

go 1.15

require (
	github.com/go-chi/chi/v5 v5.0.0
	github.com/go-chi/cors v1.2.0
	github.com/paleviews/hapi v0.1.0
	go.uber.org/zap v1.17.0
)

replace github.com/paleviews/hapi v0.1.0 => ../..
