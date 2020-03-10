package utils

import (
	"fmt"
	"net"
	"time"
)

func GetIPsFromCIDR(cidr string) ([]string, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	inc := func(ip net.IP) {
		for j := len(ip) - 1; j >= 0; j-- {
			ip[j]++
			if ip[j] > 0 {
				break
			}
		}
	}
	ips := make([]string, 0)
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	lenIPs := len(ips)
	switch {
	case lenIPs < 2:
		return ips, nil
	default:
		return ips[1 : len(ips)-1], nil
	}
}

func testConnect(m, ip string, port, timeout int) int64 {
	t := time.Now()
	addr := fmt.Sprintf("%s:%d", ip, port)
	_, err := net.DialTimeout(m, addr, time.Duration(timeout)*time.Millisecond)
	if err == nil {
		return time.Now().Sub(t).Milliseconds()
	}
	return -1
}

func TestTCP(ip string, port, timeout int) int64 {
	return testConnect("tcp",ip,port,timeout)
}

func TestUDP(ip string, port, timeout int) int64 {
	return testConnect("udp",ip,port,timeout)
}