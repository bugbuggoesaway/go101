package server

import "testing"

func TestProtocol_IsGRPCEnabled(t *testing.T) {
	tests := []struct {
		name string
		p    Protocol
		want bool
	}{
		{
			name: "grpc",
			p:    GRPC,
			want: true,
		},
		{
			name: "spex",
			p:    SPEX,
			want: false,
		},
		{
			name: "grpc|spex",
			p:    GRPC | SPEX,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.IsGRPCEnabled(); got != tt.want {
				t.Errorf("IsGRPCEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProtocol_IsSPEXEnabled(t *testing.T) {
	tests := []struct {
		name string
		p    Protocol
		want bool
	}{
		{
			name: "grpc",
			p:    GRPC,
			want: false,
		},
		{
			name: "spex",
			p:    SPEX,
			want: true,
		},
		{
			name: "grpc|spex",
			p:    GRPC | SPEX,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.IsSPEXEnabled(); got != tt.want {
				t.Errorf("IsSPEXEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}
