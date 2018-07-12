package main

import (
	"nanomsg.org/go-mangos/protocol/sub"
	"nanomsg.org/go-mangos/transport/ipc"
	"nanomsg.org/go-mangos"
	"flag"
)

func client(url string) {
	var sock mangos.Socket
	if sock, err := sub.NewSocket(); err != nil {
		println("client cannot make a socket")
		return
	}
	sock.AddTransport(ipc.NewTransport())
	if err := sock.Dial(url); err != nil {
		println("client cannot connect to server via socket")
		return
	}
	if err := sock.SetOption(mangos.OptionSubscribe, []byte("")); err != nil {
		println("client cannot subscribe topics")
		return
	}
	msg, err := sock.Recv()
	if  err != nil {
		println("client cannot receive msg via socket")
		return
	}
	println("msg received", string(msg))
}


func main() {
	url := flag.String("url", "ipc://testnanomsg", "URL of the server")
	flag.Parse()
	client(*url)
}