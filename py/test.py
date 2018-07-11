#!/usr/bin/env python3

from multiprocessing import Process
import sys
import wsclient

if __name__ == "__main__":
    clientNum = int(sys.argv[1])
    for i in range(clientNum):
        p = Process(target=wsclient.startClient, args=((i%4)+1, 2000))
        p.start()
