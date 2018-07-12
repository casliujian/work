#!/usr/bin/env python3

from multiprocessing import Process
import sys
import wsclient

if __name__ == "__main__":
    serverIp = sys.argv[1]
    serverPort = int(sys.argv[2])
    clientNum = int(sys.argv[3])
    for i in range(clientNum):
        p = Process(
            target=wsclient.startClient, 
            args=((serverIp, serverPort, i%4)+1, 2000))
        p.start()
