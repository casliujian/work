package main
import (
	"github.com/toolkits/net"
)
func main() {
	ips, _ := net.IntranetIP()
	println(ips[0])
}