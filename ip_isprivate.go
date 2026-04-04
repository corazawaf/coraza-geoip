//go:build !tinygo

package plugin

import "net"

func IsSpecialIpRangeToBeSkipped(ip net.IP) bool {
	return ip.IsPrivate() || ip.IsLoopback() || ip.IsUnspecified()
}
