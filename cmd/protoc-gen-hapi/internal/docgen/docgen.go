package docgen

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/paleviews/hapi/cmd/protoc-gen-hapi/internal/serviceregistry"
)

func Generate(reg *serviceregistry.Registry, plugin *protogen.Plugin) (string, error) {
	doc, err := transform(reg, plugin)
	if err != nil {
		return "", err
	}
	spec := doc.openAPI()
	if err := validateOpenAPI(spec); err != nil {
		return "", err
	}
	return doc.yaml()
}
