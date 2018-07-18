package main

import (
	"nanomsg.org/go-mangos/protocol/pub"
	"nanomsg.org/go-mangos/transport/ipc"
	"nanomsg.org/go-mangos/transport/tcp"
	"time"
	"flag"
)

func server(url string) {
	sock, err := pub.NewSocket()
	if err != nil {
		println("server cannot make a socket")
		return
	}
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	err = sock.Listen(url)
	if err != nil {
		println("server cannot listen on the socket")
		return
	}
	for {
		d := time.Now().Format(time.ANSIC)
		err = sock.Send([]byte(d))
		if err != nil {
			println("server cannot send msg via socket")
			return
		}
		println("sent data:", d)
		//time.Sleep(time.Microsecond)
	}

}

func main() {
	url := flag.String("url", "tcp://192.168.9.72:8000", "URL of the server")
	flag.Parse()
	println("now starting server")
	server(*url)
}