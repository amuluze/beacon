// doc-schema exports JSON field shapes from Go syntax without loading application dependencies.
package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type schema = map[string]any

func shape(e ast.Expr, prefix string) schema {
	switch t := e.(type) {
	case *ast.Ident:
		switch t.Name {
		case "string":
			return schema{"type": "string"}
		case "bool":
			return schema{"type": "boolean"}
		case "float32", "float64":
			return schema{"type": "number"}
		case "int", "int64", "int32", "uint", "uint64", "uint32", "byte":
			return schema{"type": "integer"}
		default:
			return schema{"$ref": "#/components/schemas/" + prefix + t.Name}
		}
	case *ast.ArrayType:
		if id, ok := t.Elt.(*ast.Ident); ok && id.Name == "byte" {
			return schema{"type": "string", "format": "byte"}
		}
		return schema{"type": "array", "items": shape(t.Elt, prefix), "nullable": true}
	case *ast.MapType:
		return schema{"type": "object", "additionalProperties": shape(t.Value, prefix), "nullable": true}
	case *ast.StarExpr:
		return schema{"allOf": []schema{shape(t.X, prefix)}, "nullable": true}
	case *ast.SelectorExpr:
		if id, ok := t.X.(*ast.Ident); ok && id.Name == "time" && t.Sel.Name == "Time" {
			return schema{"type": "string", "format": "date-time"}
		}
	}
	panic(fmt.Sprintf("unsupported schema expression %T", e))
}
func applyValidation(s schema, rules string) {
	for _, rule := range strings.Split(rules, ",") {
		if rule == "dive" {
			if items, ok := s["items"].(schema); ok {
				applyValidation(items, strings.SplitN(rules, "dive,", 2)[1])
			}
			break
		}
		if rule == "required" {
			delete(s, "nullable")
			// validator required rejects scalar zero values, not merely absent keys.
			// Non-nil empty slices/maps and pointers to zero values remain valid.
			switch s["type"] {
			case "string":
				s["minLength"] = 1
			case "integer", "number":
				s["not"] = schema{"enum": []any{0}}
			case "boolean":
				s["enum"] = []any{true}
			}
		}
		if rule == "uuid" {
			s["format"] = "uuid"
		}
		if rule == "email" {
			s["format"] = "email"
		}
		pair := strings.SplitN(rule, "=", 2)
		if len(pair) != 2 {
			continue
		}
		if pair[0] == "oneof" {
			values := []any{}
			for _, v := range strings.Fields(pair[1]) {
				if s["type"] == "integer" {
					n, err := strconv.Atoi(v)
					if err != nil {
						panic(err)
					}
					values = append(values, n)
				} else {
					values = append(values, v)
				}
			}
			s["enum"] = values
			continue
		}
		n, err := strconv.Atoi(pair[1])
		if err != nil {
			continue
		}
		low, high := "minimum", "maximum"
		switch s["type"] {
		case "string":
			low, high = "minLength", "maxLength"
		case "array":
			low, high = "minItems", "maxItems"
		case "object":
			low, high = "minProperties", "maxProperties"
		}
		switch pair[0] {
		case "min", "gte":
			s[low] = n
		case "max", "lte":
			s[high] = n
		case "len":
			s[low] = n
			s[high] = n
		case "gt":
			s[low] = n
			s["exclusiveMinimum"] = true
		}
	}
}
func main() {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(exportSchemas(repoRoot())); err != nil {
		panic(err)
	}
}

// repoRoot locates the workspace root through go.work so the exported shapes
// never depend on the caller's working directory. Isolated fixtures without a
// workspace file fall back to the working directory.
func repoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	start := dir
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	for _, required := range []string{"beacon/service/schema", "common/rpc/schema"} {
		if _, err := os.Stat(filepath.Join(start, required)); err != nil {
			panic("go.work not found and " + required + " is missing from " + start)
		}
	}
	return start
}

// exportSchemas derives the JSON field shapes from the Go schema sources.
func exportSchemas(root string) map[string]schema {
	out := map[string]schema{}
	for _, group := range []struct {
		pattern, prefix string
		selected        map[string]string
	}{
		{filepath.Join(root, "beacon/service/schema/*.go"), "", nil},
		{filepath.Join(root, "common/rpc/schema/report.go"), "Report", nil},
		{filepath.Join(root, "beacon/service/agent/updater.go"), "", map[string]string{"UpdateResult": "UpgradeResult"}},
		{filepath.Join(root, "beacon/service/agent/agent.go"), "", map[string]string{"upgradeArgs": "UpgradeArgs"}},
		{filepath.Join(root, "beacon/service/model/agent.go"), "", map[string]string{"Agent": "ModelAgent"}},
		{filepath.Join(root, "beacon/pkg/fiberx/fiberx.go"), "", map[string]string{"FailedResponse": "ErrorResponse"}},
	} {
		paths, err := filepath.Glob(group.pattern)
		if err != nil {
			panic(err)
		}
		if len(paths) == 0 {
			panic("missing schema source: " + group.pattern)
		}
		for _, path := range paths {
			if strings.HasSuffix(path, "_test.go") {
				continue
			}
			f, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
			if err != nil {
				panic(err)
			}
			for _, decl := range f.Decls {
				gen, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}
				for _, spec := range gen.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					if group.selected != nil && group.selected[ts.Name.Name] == "" {
						continue
					}
					st, ok := ts.Type.(*ast.StructType)
					if !ok {
						continue
					}
					props := map[string]schema{}
					required := []string{}
					embedded := []schema{}
					for _, field := range st.Fields.List {
						if len(field.Names) == 0 && ts.Name.Name == "Agent" {
							if sel, ok := field.Type.(*ast.SelectorExpr); ok {
								if id, ok := sel.X.(*ast.Ident); ok && id.Name == "gorm" && sel.Sel.Name == "Model" {
									embedded = append(embedded, schema{"$ref": "#/components/schemas/GORMModel"})
									continue
								}
							}
						}
						if len(field.Names) != 1 {
							panic("embedded or grouped field requires explicit modeling")
						}
						name := field.Names[0].Name
						if !ast.IsExported(name) {
							continue
						}
						tag := reflect.StructTag("")
						if field.Tag != nil {
							raw, err := strconv.Unquote(field.Tag.Value)
							if err != nil {
								panic(err)
							}
							tag = reflect.StructTag(raw)
						}
						parts := strings.Split(tag.Get("json"), ",")
						if parts[0] == "-" {
							continue
						}
						if parts[0] != "" {
							name = parts[0]
						}
						p := shape(field.Type, group.prefix)
						if description := tag.Get("description"); description != "" {
							p["description"] = description
						}
						if pattern := tag.Get("pattern"); pattern != "" {
							p["pattern"] = pattern
						}
						if query := tag.Get("query"); query != "" {
							p["x-go-query"] = strings.Split(query, ",")[0]
						}
						// Preserve every rule and translate basic bounds/enums. Custom validators
						// and cross-field business rules still require handler/service tests.
						if v := tag.Get("validate"); v != "" {
							p["x-go-validate"] = v
							applyValidation(p, v)
						}
						props[name] = p
						isArgs := strings.HasSuffix(ts.Name.Name, "Args")
						if (!isArgs && !strings.Contains(tag.Get("json"), "omitempty")) || (isArgs && strings.Contains(","+tag.Get("validate")+",", ",required,")) {
							required = append(required, name)
						}
					}
					source := path
					if rel, err := filepath.Rel(root, path); err == nil {
						source = filepath.ToSlash(rel)
					}
					s := schema{"type": "object", "properties": props, "x-go-source": source + ":" + ts.Name.Name}
					// JSON output fields are required; input requirements are enforced by Go validators.
					if len(required) > 0 {
						s["required"] = required
					}
					if len(embedded) > 0 {
						s["allOf"] = embedded
					}
					name := group.prefix + ts.Name.Name
					if group.selected != nil {
						name = group.selected[ts.Name.Name]
					}
					if name == "UpgradeArgs" {
						// Service validation requires an explicit target.
						s["required"] = []string{"agent_id"}
					}
					out[name] = s
				}
			}
		}
	}
	return out
}
