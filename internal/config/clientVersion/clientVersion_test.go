package clientVersion

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatClientVersion(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
	}{
		{
			name:    "test_format_1",
			args:    args{id: "2.0.1"},
			want:    32769,
			wantErr: false,
		},
		{
			name:    "test_format_max",
			args:    args{id: "3.3.4095"},
			want:    65535,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatClientVersion(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatClientVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FormatClientVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseClientVersion(t *testing.T) {
	type args struct {
		id uint16
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test_parse_1",
			args:    args{id: 32769},
			want:    "2.0.1",
			wantErr: false,
		},
		{
			name:    "test_format_max",
			args:    args{id: 65535},
			want:    "3.3.4095",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseClientVersion(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseClientVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseClientVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatClientVersion_2(t *testing.T) {
	id, err := FormatClientVersion("1.0.3")
	assert.NoError(t, err)
	t.Log(id)
}
