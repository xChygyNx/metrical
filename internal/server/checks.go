package server

import (
	"fmt"
	"net"

	"github.com/xChygyNx/metrical/internal/server/types"
)

func isIpInSubnet(addr string, conf *types.SyncInfo) (res bool, err error) {
	cidr := conf.TrustedSubnet

	if cidr != "" {
		ip := net.ParseIP(addr)
		if ip == nil {
			return false, fmt.Errorf("error in parse IP address: %w", err)
		}

		_, subnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return false, fmt.Errorf("error in parse CIDR: %w", err)
		}

		res = subnet.Contains(ip)
		return res, nil
	}
	return true, nil
}
