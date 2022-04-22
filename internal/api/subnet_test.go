package api

import (
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestSubnet_readMasks(t *testing.T) {
	tests := []struct {
		name   string
		cfg    string
		expect []net.IPNet
	}{
		{name: "empty", cfg: "", expect: []net.IPNet{}},
		{name: "all ok", cfg: "172.17.0.0/16, 10.10.20.0/30",
			expect: []net.IPNet{
				{IP: net.IPv4(172, 17, 0, 0), Mask: net.CIDRMask(16, 32)},
				{IP: net.IPv4(10, 10, 20, 0), Mask: net.CIDRMask(30, 32)},
			},
		},
		{name: "second wrong", cfg: "172.17.0.0/16, 300.10.20.0/30, 10.10.20.0/16",
			expect: []net.IPNet{
				{IP: net.IPv4(172, 17, 0, 0), Mask: net.CIDRMask(16, 32)},
				{IP: net.IPv4(10, 10, 0, 0), Mask: net.CIDRMask(16, 32)},
			},
		},
		{name: "all wrong", cfg: "172 17.0.1/16, 300.10.20.0/30, 10.10.20.0/60",
			expect: []net.IPNet{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			masks := readMasks(tt.cfg)
			require.Equal(t, len(tt.expect), len(masks), "count of masks not expected")
			for i, v := range masks {
				//  compare by strings CIDRMask parse returns 16 bit in ip, net.IPv4 4 bit
				require.Equal(t, tt.expect[i].IP.String(), v.IP.String())
				require.Equal(t, tt.expect[i].Mask, v.Mask)
			}
		})
	}
}
