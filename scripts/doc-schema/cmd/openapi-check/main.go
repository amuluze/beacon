// openapi-check validates an OpenAPI 3 document supplied as JSON on stdin.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

func main() {
	schemaName := flag.String("schema", "", "component schema to validate a sample against")
	sample := flag.String("sample", "", "JSON sample for the component schema")
	flag.Parse()
	data, err := io.ReadAll(os.Stdin)
	if err == nil {
		var doc *openapi3.T
		loader := openapi3.NewLoader() // External references are disabled.
		doc, err = loader.LoadFromData(data)
		if err == nil {
			err = doc.Validate(context.Background())
		}
		if err == nil {
			err = validateSecurity(doc)
		}
		if err == nil && *schemaName != "" {
			var value any
			err = json.Unmarshal([]byte(*sample), &value)
			if err == nil {
				if doc.Components == nil || doc.Components.Schemas[*schemaName] == nil || doc.Components.Schemas[*schemaName].Value == nil {
					err = fmt.Errorf("undefined component schema %q", *schemaName)
				} else {
					err = doc.Components.Schemas[*schemaName].Value.VisitJSON(value)
				}
			}
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("OpenAPI 3 standard validation passed")
}

// kin-openapi validates scheme objects but does not resolve security requirement
// names. Check those cross-references for both root and operation overrides.
func validateSecurity(doc *openapi3.T) error {
	check := func(requirements openapi3.SecurityRequirements) error {
		for _, requirement := range requirements {
			for name, scopes := range requirement {
				if doc.Components == nil || doc.Components.SecuritySchemes[name] == nil || doc.Components.SecuritySchemes[name].Value == nil {
					return fmt.Errorf("undefined security scheme %q", name)
				}
				scheme := doc.Components.SecuritySchemes[name].Value
				if scheme.Type != "oauth2" && scheme.Type != "openIdConnect" && len(scopes) > 0 {
					return fmt.Errorf("security scheme %q does not support scopes", name)
				}
			}
		}
		return nil
	}
	if err := check(doc.Security); err != nil {
		return err
	}
	for path, item := range doc.Paths.Map() {
		for method, operation := range item.Operations() {
			if operation.Security != nil {
				if err := check(*operation.Security); err != nil {
					return fmt.Errorf("%s %s: %w", method, path, err)
				}
			}
		}
	}
	return nil
}
