package main

import "net"

func isHostname(s string) bool {
	if _, err := net.LookupHost(s); err != nil {
		return false
	}
	return true
}

func isIP(s string) bool {
	return net.ParseIP(s) != nil
}
