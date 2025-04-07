package server

import (
	"net"
	"testing"
)

func Test_broadcastFromCIDR(t *testing.T) {
	tests := []struct {
		name    string
		addr    net.Addr
		want    net.IP
		wantErr bool
	}{
		{
			name: "/32 address",
			addr: &net.IPNet{
				IP:   net.IPv4(1, 1, 1, 1),
				Mask: net.IPMask{255, 255, 255, 255},
			},
			want: net.IPv4(1, 1, 1, 1),
		},
		{
			name: "/24 address",
			addr: &net.IPNet{
				IP:   net.IPv4(1, 1, 1, 1),
				Mask: net.IPMask{255, 255, 255, 0},
			},
			want: net.IPv4(1, 1, 1, 255),
		},
		{
			name: "invalid ip",
			addr: &net.IPNet{
				IP: net.IP{1, 1, 1},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ipv6",
			addr: &net.IPNet{
				IP:   net.IPv6loopback,
				Mask: net.CIDRMask(64, 128),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := broadcastFromCIDR(tt.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("broadcastFromCIRD() err = %v, wantErr %v", err, tt.wantErr)
			}
			if !got.Equal(tt.want) {
				t.Errorf("broadcastFromCIDR() = %v, want %v", got, tt.want)
			}
		})
	}
}
