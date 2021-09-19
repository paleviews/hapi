package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/paleviews/hapi/cmd/protoc-gen-hapi/internal/codegen"
	"github.com/paleviews/hapi/cmd/protoc-gen-hapi/internal/docgen"
	"github.com/paleviews/hapi/cmd/protoc-gen-hapi/internal/serviceregistry"
)

const version = "v0.1.0"

func main() {
	showVersion := flag.Bool("version", false, "show version")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-hapi %v\n", version)
		return
	}

	var flags flag.FlagSet
	docPath := flags.String("doc_file", "", "file path of doc")
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		if *docPath == "" {
			err := errors.New("no doc file provided")
			gen.Error(err)
			return err
		}

		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		reg, err := serviceregistry.Load(gen)
		if err != nil {
			gen.Error(err)
			return err
		}

		docContent, err := docgen.Generate(reg, gen)
		if err != nil {
			gen.Error(err)
			return err
		}
		err = os.WriteFile(*docPath, []byte(docContent), 0644)
		if err != nil {
			gen.Error(err)
			return err
		}

		codegen.Generate(reg, gen, version)
		return nil
	})
}
