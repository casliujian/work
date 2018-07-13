package main

import (
	"flag"
	"nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/sub"
	"nanomsg.org/go-mangos/transport/ipc"
	"nanomsg.org/go-mangos/transport/tcp"
)

func client(url string) {
	sock, err := sub.NewSocket()
	if err != nil {
		println("client cannot make a socket")
		return
	}
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Dial(url); err != nil {
		println("client cannot connect to server via socket")
		return
	}
	if err = sock.SetOption(mangos.OptionSubscribe, []byte("")); err != nil {
		println("client cannot subscribe topics")
		return
	}
	msg := make([]byte, 4)
	for {
		msg, err = sock.Recv()
		println("size of msg:", len(msg))
		if err != nil {
			println("client cannot receive msg via socket")
			return
		}
		println("msg received", string(msg))
	}
}

func main() {
	url := flag.String("url", "tcp://192.168.9.72:8000", "URL of the server")
	flag.Parse()
	println("now starting client")
	client(*url)
}
