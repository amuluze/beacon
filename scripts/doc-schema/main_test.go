package main

import (
	"encoding/json"
	"go/parser"
	"reflect"
	"testing"
)

func shapeOf(t *testing.T, expr string) schema {
	t.Helper()
	e, err := parser.ParseExpr(expr)
	if err != nil {
		t.Fatalf("parse %q: %v", expr, err)
	}
	return shape(e, "")
}

// assertSchema compares a schema against its JSON form so key order never matters.
func assertSchema(t *testing.T, got schema, want string) {
	t.Helper()
	var expected map[string]any
	if err := json.Unmarshal([]byte(want), &expected); err != nil {
		t.Fatalf("bad expectation %q: %v", want, err)
	}
	encoded, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var normalized map[string]any
	if err := json.Unmarshal(encoded, &normalized); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(normalized, expected) {
		t.Errorf("got %s, want %s", encoded, want)
	}
}

func TestShapeScalarsAndComposites(t *testing.T) {
	cases := []struct {
		typ  string
		want string
	}{
		{"string", `{"type":"string"}`},
		{"bool", `{"type":"boolean"}`},
		{"int", `{"type":"integer"}`},
		{"float64", `{"type":"number"}`},
		{"Other", `{"$ref":"#/components/schemas/Other"}`},
		{"time.Time", `{"type":"string","format":"date-time"}`},
		{"[]byte", `{"type":"string","format":"byte"}`},
		{"[]string", `{"type":"array","items":{"type":"string"},"nullable":true}`},
		{"map[string]int", `{"type":"object","additionalProperties":{"type":"integer"},"nullable":true}`},
		{"*int", `{"allOf":[{"type":"integer"}],"nullable":true}`},
	}
	for _, c := range cases {
		t.Run(c.typ, func(t *testing.T) {
			assertSchema(t, shapeOf(t, c.typ), c.want)
		})
	}
}

func TestShapeUsesSchemaPrefix(t *testing.T) {
	e, err := parser.ParseExpr("ReportItem")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	assertSchema(t, shape(e, "Report"), `{"$ref":"#/components/schemas/ReportReportItem"}`)
}

func TestApplyValidationRules(t *testing.T) {
	cases := []struct {
		name  string
		in    schema
		rules string
		want  string
	}{
		{"required drops nullable", schema{"type": "string", "nullable": true}, "required", `{"type":"string","minLength":1}`},
		{"uuid sets format", schema{"type": "string"}, "uuid", `{"type":"string","format":"uuid"}`},
		{"email sets format", schema{"type": "string"}, "email", `{"type":"string","format":"email"}`},
		{"string bounds", schema{"type": "string"}, "min=1,max=128", `{"type":"string","minLength":1,"maxLength":128}`},
		{"integer bounds", schema{"type": "integer"}, "gte=1,lte=10", `{"type":"integer","minimum":1,"maximum":10}`},
		{"oneof enum", schema{"type": "string"}, "oneof=a b c", `{"type":"string","enum":["a","b","c"]}`},
		{"exclusive minimum", schema{"type": "integer"}, "gt=0", `{"type":"integer","minimum":0,"exclusiveMinimum":true}`},
		{"array item bounds", schema{"type": "array", "items": schema{"type": "string", "nullable": true}}, "dive,min=2", `{"type":"array","items":{"type":"string","nullable":true,"minLength":2}}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			applyValidation(c.in, c.rules)
			assertSchema(t, c.in, c.want)
		})
	}
}

func TestExportSchemasCoversWorkspaceSources(t *testing.T) {
	out := exportSchemas(repoRoot())
	if len(out) < 100 {
		t.Fatalf("expected the schema inventory to cover the workspace, got %d entries", len(out))
	}
	for _, name := range []string{"UserDeleteArgs", "LoginArgs"} {
		s, ok := out[name]
		if !ok {
			t.Fatalf("missing exported schema %s", name)
		}
		if s["type"] != "object" {
			t.Errorf("%s type = %v, want object", name, s["type"])
		}
		if _, ok := s["properties"]; !ok {
			t.Errorf("%s has no properties", name)
		}
	}
	for _, name := range []string{"CPUUsageArgs", "MemoryUsageArgs", "DiskUsageArgs", "NetworkUsageArgs", "ContainerUsageArgs"} {
		props := out[name]["properties"].(map[string]schema)
		for field, wire := range map[string]string{"StartTime": "start_time", "EndTime": "end_time"} {
			if got := props[field]["x-go-query"]; got != wire {
				t.Errorf("%s.%s query name = %v, want %s", name, field, got, wire)
			}
		}
		if _, renamed := props["start_time"]; renamed {
			t.Errorf("query alias must not rename JSON field in %s", name)
		}
	}
}

func TestRequiredZeroValueSemantics(t *testing.T) {
	for _, tc := range []struct{ typ, want string }{
		{"int", `{"type":"integer","not":{"enum":[0]}}`},
		{"bool", `{"type":"boolean","enum":[true]}`},
		{"[]string", `{"type":"array","items":{"type":"string"}}`},
		{"*int", `{"allOf":[{"type":"integer"}]}`},
	} {
		t.Run(tc.typ, func(t *testing.T) {
			s := shapeOf(t, tc.typ)
			applyValidation(s, "required")
			assertSchema(t, s, tc.want)
		})
	}
}
