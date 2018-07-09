#!/usr/bin/env python3

import socket
import array
import sys

port = 2000
serversocket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
serversocket.bind(('',port))
serversocket.listen()
channel1 = b'\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01'
channel2 = b'\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02'
channel3 = b'\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x03'
channel4 = b'\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x04'

doublearray = array.array('d')
for i in range(65536):
    doublearray.insert(i,0.1+i/65536)
databytes = doublearray.tobytes()

clientsocket, addr = serversocket.accept()

def exitfunc():
    print('data source exit')
    clientsocket.close()
    serversocket.close()
# sys.execfunc = exitfunc
try:
    while True:
        clientsocket.send(channel1)
        clientsocket.send(databytes)
        clientsocket.send(channel2)
        clientsocket.send(databytes)
        clientsocket.send(channel3)
        clientsocket.send(databytes)
        clientsocket.send(channel4)
        clientsocket.send(databytes)
except KeyboardInterrupt:
    exitfunc()


