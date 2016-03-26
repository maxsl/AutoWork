//安全起见,设置允许的网段,只允许内网.
package fileagent

import "net"

var (
	mask  net.IPMask = net.CIDRMask(16, 32)
	IP192 net.IPNet  = net.IPNet{net.ParseIP("192.168.0.0"), mask}
	IP172 net.IPNet  = net.IPNet{net.ParseIP("172.18.0.0"), mask}
	IP10  net.IPNet  = net.IPNet{net.ParseIP("10.0.0.0"), mask}
)

const (
	ForbiddenAccess = 403
	MarshalError    = 550
	UnmarshalError  = 551
	ReadBodyError   = 552
)

func ipIsLanIP(str string) bool {
	ip := net.ParseIP(str)
	return IP192.Contains(ip) || IP172.Contains(ip) || IP10.Contains(ip)
}
