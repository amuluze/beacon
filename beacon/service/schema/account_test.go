package schema

import (
	"fmt"
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

func TestPasswordUpdateArgsRejectsShortNewPassword(t *testing.T) {
	for _, tc := range []struct {
		password string
		invalid  bool
	}{
		{"x", true}, {"abcDEF123456", false}, {"密码abcDEF123456", false},
	} {
		err := validatex.ValidateStruct(&PasswordUpdateArgs{Username: "admin", OldPassword: "old-password", NewPassword: tc.password})
		if (err != nil) != tc.invalid {
			t.Fatalf("password length %d: unexpected validation result", len([]rune(tc.password)))
		}
	}
}

// Keep these samples paired with the OpenAPI gate's business samples.
// Create status is required; zero is not an account state.
func TestRequiredScalarContractSamples(t *testing.T) {
	for _, tc := range []struct {
		name     string
		args     any
		accepted bool
	}{
		{"empty path", &FilesSearchArgs{Path: ""}, false},
		{"nonempty path", &FilesSearchArgs{Path: "/"}, true},
		{"zero status", &UserCreateArgs{Username: "audituser", Password: "Strong-Pass9!", Status: 0, RoleIDs: []string{"r"}}, false},
		{"active status", &UserCreateArgs{Username: "audituser", Password: "Strong-Pass9!", Status: 1, RoleIDs: []string{"r"}}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := validatex.ValidateStruct(tc.args); (err == nil) != tc.accepted {
				t.Fatalf("unexpected acceptance: %v", err)
			}
		})
	}
}

func TestAccountStatusValidation(t *testing.T) {
	for _, status := range []int{-1, 0, 1, 2, 3} {
		for _, tc := range []struct {
			name  string
			args  any
			valid bool
		}{
			{"create user", &UserCreateArgs{Username: "audituser", Password: "Strong-Pass9!", Status: status, RoleIDs: []string{"r"}}, status == 1 || status == 2},
			{"create role", &RoleCreateArgs{Name: "operator", Status: status}, status == 1 || status == 2},
			{"update user", &UserUpdateArgs{ID: "019f8dbc-8781-7af3-bf49-e7fafa912e87", Status: status}, status >= 0 && status <= 2},
			{"update role", &RoleUpdateArgs{ID: "019f8dbc-8781-7af3-bf49-e7fafa912e87", Status: status}, status >= 0 && status <= 2},
		} {
			t.Run(tc.name+fmt.Sprint(status), func(t *testing.T) {
				if err := validatex.ValidateStruct(tc.args); (err == nil) != tc.valid {
					t.Fatalf("status %d: %v", status, err)
				}
			})
		}
	}
}
