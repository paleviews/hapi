package docgen

import (
	"strings"
	"testing"

	"github.com/paleviews/hapi/cmd/protoc-gen-hapi/internal/serviceregistry"
)

func TestDocumentYAMLRetainsCommentsAndBlockDescriptions(t *testing.T) {
	doc := &document{
		info: docInfo{
			version: "v0.0.1",
			title:   "test doc",
		},
		endpoints: []*endpoint{
			{
				path: "/v1/{id}",
				operations: []*operation{
					{
						httpMethod:  serviceregistry.HTTPMethodGet,
						description: description{"operation line 1", "operation line 2"},
						operationID: "pkg.Service.Get",
						parameters: []parameter{
							&pathParameter{
								name:       "id",
								scalarType: builtinString,
							},
						},
						response: anonymousObject{
							{
								name:        "code",
								description: description{"response code enums:", ".. 0: ok"},
								_type:       builtinInt32,
								possibilities: possibilities{
									{literal: "0", comment: "ok"},
								},
							},
							{
								name:  "message",
								_type: builtinString,
								possibilities: possibilities{
									{literal: "'ok'", comment: "0"},
								},
							},
							{
								name:        "direction",
								description: description{"direction line 1", "direction line 2"},
								_type: enum{
									{name: "DIRECTION_UNKNOWN", description: description{"one line enum"}},
									{name: "EAST"},
									{name: "SOUTH", description: description{"multiple lines 1", "multiple lines 2"}},
								},
							},
						},
					},
				},
			},
		},
	}

	got := renderTestDocument(t, doc)
	assertContainsAll(t, got,
		`openapi: "3.0.3"`,
		"      description: |\n        operation line 1\\\n        operation line 2\n",
		"                    description: |\n                      response code enums:\\\n                      .. 0: ok\n",
		"                      - 0 # ok\n",
		"                      - 'ok' # 0\n",
		"                      - DIRECTION_UNKNOWN # one line enum\n",
		"                      # multiple lines 1\n                      # multiple lines 2\n                      - SOUTH\n",
	)
}

func TestDocumentYAMLKeepsRefsAndInlineSchemas(t *testing.T) {
	inline := &schema{
		id:       "pkg.Inline",
		refCount: 1,
		obj: anonymousObject{
			{name: "name", _type: builtinString},
		},
	}
	reused := &schema{
		id:       "pkg.Reused",
		refCount: 2,
		obj: anonymousObject{
			{name: "id", _type: builtinString},
		},
	}
	doc := &document{
		info: docInfo{
			version: "v0.0.1",
			title:   "test doc",
		},
		endpoints: []*endpoint{
			{
				path: "/v1/{id}",
				operations: []*operation{
					{
						httpMethod:  serviceregistry.HTTPMethodGet,
						operationID: "pkg.Service.Get",
						parameters: []parameter{
							&pathParameter{
								name:       "id",
								scalarType: builtinString,
							},
						},
						response: anonymousObject{
							{name: "inline", _type: inline},
							{name: "reused", _type: reused},
						},
					},
				},
			},
		},
		schemas: schemas{inline, reused},
	}

	got := renderTestDocument(t, doc)
	assertContainsAll(t, got,
		"                  inline:\n                    type: object\n                    properties:\n                      name:\n                        type: string\n",
		"                  reused:\n                    $ref: '#/components/schemas/pkg.Reused'\n",
		"  schemas:\n    pkg.Reused:\n      type: object\n",
	)
	if strings.Contains(got, "pkg.Inline:") {
		t.Fatalf("inline schema should not be emitted as a component:\n%s", got)
	}
}

func renderTestDocument(t *testing.T, doc *document) string {
	t.Helper()
	if err := validateOpenAPI(doc.openAPI()); err != nil {
		t.Fatalf("OpenAPI document did not validate: %v", err)
	}
	got, err := doc.yaml()
	if err != nil {
		t.Fatalf("rendering YAML: %v", err)
	}
	return got
}

func assertContainsAll(t *testing.T, haystack string, needles ...string) {
	t.Helper()
	for _, needle := range needles {
		if !strings.Contains(haystack, needle) {
			t.Fatalf("expected YAML to contain %q:\n%s", needle, haystack)
		}
	}
}
