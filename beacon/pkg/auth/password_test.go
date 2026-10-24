package auth

import "testing"

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{name: "strong", password: "Correct-Horse-9"},
		{name: "strong passphrase", password: "beacon secure phrase 9"},
		{name: "too short", password: "Short-9", wantErr: true},
		{name: "too few classes", password: "onlylowercase", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
