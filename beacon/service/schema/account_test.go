package schema

import (
	"testing"

	"beacon/pkg/validatex"
)

func TestUserDeleteArgsRequiresExplicitUUIDTargets(t *testing.T) {
	tests := []struct {
		name    string
		args    UserDeleteArgs
		wantErr bool
	}{
		{name: "valid", args: UserDeleteArgs{IDs: []string{"019f8dbc-8781-7af3-bf49-e7fafa912e87"}}},
		{name: "nil", args: UserDeleteArgs{}, wantErr: true},
		{name: "empty", args: UserDeleteArgs{IDs: []string{}}, wantErr: true},
		{name: "empty element", args: UserDeleteArgs{IDs: []string{""}}, wantErr: true},
		{name: "invalid length", args: UserDeleteArgs{IDs: []string{"user-1"}}, wantErr: true},
		{name: "invalid UUID", args: UserDeleteArgs{IDs: []string{"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatex.ValidateStruct(&tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
