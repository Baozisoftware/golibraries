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

func testConnect(m, ip string, port, timeout, count int) (min, max, avg, out int) {
	min = -1
	max = -1
	avg = -1
	for count < 1 {
		count = 1
	}
	sum := 0
	c := 0
	for i := 1; i <= count; i++ {
		t := time.Now()
		addr := fmt.Sprintf("%s:%d", ip, port)
		conn, err := net.DialTimeout(m, addr, time.Duration(timeout)*time.Millisecond)
		if err == nil {
			ms := int(time.Now().Sub(t).Milliseconds())
			sum += ms
			c++
			if min == -1 {
				min = ms
			} else if ms < min {
				min = ms
			}
			if max == -1 {
				max = ms
			} else if ms > max {
				max = ms
			}
		} else {
			out++
		}
		if conn != nil {
			_ = conn.Close()
		}
	}
	if c > 0 {
		avg = sum / c
	}
	return
}

func TestTCPConnect(ip string, port, timeout, count int) (min, max, avg, out int) {
	return testConnect("tcp", ip, port, timeout, count)
}

func TestUDPConnect(ip string, port, timeout, count int) (min, max, avg, out int) {
	return testConnect("udp", ip, port, timeout, count)
}
