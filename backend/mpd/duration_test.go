package mpd

import (
	"testing"
	"time"
)

func TestISO8601Duration_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		want    ISO8601Duration
		wantErr bool
	}{
		{
			name:    "integer",
			text:    "PT1S",
			want:    ISO8601Duration(time.Second),
			wantErr: false,
		},
		{
			name:    "big integer",
			text:    "PT18378S",
			want:    ISO8601Duration(time.Second * 18378),
			wantErr: false,
		},
		{
			name:    "decimal",
			text:    "PT49.006S",
			want:    ISO8601Duration(49.006 * float64(time.Second)),
			wantErr: false,
		},
		{
			name:    "big decimal",
			text:    "PT18378.5S",
			want:    ISO8601Duration(18378.5 * float64(time.Second)),
			wantErr: false,
		},
		{
			name:    "very big decimal",
			text:    "PT1837832819.5321321S",
			want:    ISO8601Duration(1837832819.5321321 * float64(time.Second)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dummy ISO8601Duration
			if err := dummy.UnmarshalText([]byte(tt.text)); (err != nil) != tt.wantErr {
				t.Errorf("ISO8601Duration.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.want != dummy {
					t.Errorf("ISO8601Duration.UnmarshalText() = %v, want %v", dummy, tt.want)
				}
			}
		})
	}
}

func TestISO8601Duration_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		i       ISO8601Duration
		want    string
		wantErr bool
	}{
		{
			name:    "integer",
			i:       ISO8601Duration(time.Second),
			want:    "PT1S",
			wantErr: false,
		},
		{
			name:    "big integer",
			i:       ISO8601Duration(18378 * time.Second),
			want:    "PT18378S",
			wantErr: false,
		},
		{
			name:    "decimal",
			i:       ISO8601Duration(49.006 * float64(time.Second)),
			want:    "PT49.006S",
			wantErr: false,
		},
		{
			name:    "big decimal",
			i:       ISO8601Duration(18378.5 * float64(time.Second)),
			want:    "PT18378.5S",
			wantErr: false,
		},
		{
			name:    "very big decimal",
			i:       ISO8601Duration(1837832819.5321321 * float64(time.Second)),
			want:    "PT1837832819.5321321S",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("ISO8601Duration.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("ISO8601Duration.MarshalText() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
