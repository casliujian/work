package main

import (
	"nanomsg.org/go-mangos/protocol/pub"
	"nanomsg.org/go-mangos/transport/ipc"
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
	err := sock.Listen(url)
	if err != nil {
		println("server cannot listen on the socket")
		return
	}
	d := time.Now().Format(time.ANSIC)
	err := sock.Send([]byte(d))
	if err != nil {
		println("server cannot send msg via socket")
		return
	}
}

func main() {
	url := flag.String("url", "ipc://testnanomsg", "URL of the server")
	flag.Parse()
	server(*url)
}