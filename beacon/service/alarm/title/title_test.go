package title

import (
	"strings"
	"testing"

	"beacon/service/model"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{name: "trim", raw: "  自定义告警标题  ", want: "自定义告警标题"},
		{name: "empty", raw: "   ", wantErr: true},
		{name: "line feed", raw: "告警\n标题", wantErr: true},
		{name: "carriage return", raw: "告警\r标题", wantErr: true},
		{name: "maximum length", raw: strings.Repeat("告", MaxLength), want: strings.Repeat("告", MaxLength)},
		{name: "too long", raw: strings.Repeat("告", MaxLength+1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Normalize(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Normalize() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("Normalize() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestOrDefault(t *testing.T) {
	if got := OrDefault("\n"); got != model.DefaultAlarmMessageTitle {
		t.Fatalf("OrDefault() = %q, want %q", got, model.DefaultAlarmMessageTitle)
	}
	if got := OrDefault(" 自定义标题 "); got != "自定义标题" {
		t.Fatalf("OrDefault() = %q, want normalized custom title", got)
	}
}
