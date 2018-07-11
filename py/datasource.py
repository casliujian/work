#!/usr/bin/env python3

import socket
import array
import sys

port = 2000
serversocket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
serversocket.bind(('',port))
serversocket.listen()
channel1 = b'\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
channel2 = b'\x02\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
channel3 = b'\x03\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
channel4 = b'\x04\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'

doublearray = array.array('d')
for i in range(65536):
    doublearray.insert(i,0.000001+i)
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


