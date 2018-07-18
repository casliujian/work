package main

import "github.com/toolkits/net"

var localIP = ""

func getLocalIP() string {
	ips, _ := net.IntranetIP()
	return ips[0]
}
