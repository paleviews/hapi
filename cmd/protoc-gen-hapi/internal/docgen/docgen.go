package docgen

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/paleviews/hapi/cmd/protoc-gen-hapi/internal/printer"
	"github.com/paleviews/hapi/cmd/protoc-gen-hapi/internal/serviceregistry"
)

func Generate(reg *serviceregistry.Registry, plugin *protogen.Plugin) (string, error) {
	doc, err := transform(reg, plugin)
	if err != nil {
		return "", err
	}
	p := printer.New("  ")
	doc.print(p)
	return p.Content(), nil
}
