//go:build tinygo

package plugin

import "net"

func IsSpecialIpRangeToBeSkipped(ip net.IP) bool {
	//"IsPrivate" is not defined for tinygo, see: tinygo.org/docs/reference/lang-support/stdlib/#netnetip
	// therfore we return false here
	return false
}
