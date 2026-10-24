package schema

import (
	"strings"
	"testing"

	"beacon/pkg/validatex"
)

func TestAlarmThresholdUpdateArgsRequiresCanonicalType(t *testing.T) {
	tests := []struct {
		name    string
		args    AlarmThresholdUpdateArgs
		wantErr bool
	}{
		{name: "cpu", args: AlarmThresholdUpdateArgs{ID: 1, Type: "cpu", Duration: 0, Threshold: 0}},
		{name: "memory", args: AlarmThresholdUpdateArgs{ID: 2, Type: "memory", Duration: 2, Threshold: 80}},
		{name: "disk", args: AlarmThresholdUpdateArgs{ID: 3, Type: "disk", Duration: 2, Threshold: 85}},
		{name: "legacy memory alias", args: AlarmThresholdUpdateArgs{ID: 2, Type: "mem", Duration: 2, Threshold: 80}, wantErr: true},
		{name: "missing type", args: AlarmThresholdUpdateArgs{ID: 1, Duration: 2, Threshold: 80}, wantErr: true},
		{name: "percentage above range", args: AlarmThresholdUpdateArgs{ID: 1, Type: "cpu", Duration: 2, Threshold: 101}, wantErr: true},
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

func TestAlarmMessageSettingUpdateArgsValidation(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{name: "valid", title: "Beacon 主机告警"},
		{name: "maximum length", title: strings.Repeat("告", 128)},
		{name: "empty", title: "", wantErr: true},
		{name: "too long", title: strings.Repeat("告", 129), wantErr: true},
		// CR/LF and whitespace-only values require semantic validation in the
		// service after trimming; schema validation still exercises the API bound.
		{name: "whitespace reaches service", title: "   "},
		{name: "line feed reaches service", title: "告警\n标题"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := AlarmMessageSettingUpdateArgs{MessageTitle: tt.title}
			err := validatex.ValidateStruct(&args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
