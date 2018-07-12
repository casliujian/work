#!/usr/bin/env python3

import asyncio
import websockets
import time
import sys
import array
import struct
import json

# recvno = 40000


async def hello(uri, pipeNum, dataItemNum):
    recvd = 0
    async with websockets.connect(uri) as websocket:
        jsonData = json.dumps({'pipeNum': pipeNum, 'dataItemNum': dataItemNum})
        await websocket.send(jsonData)
        print("json data sent")

        while True:
            msg = await websocket.recv()
            da = array.array('d')
            da.frombytes(msg)
            print("received", da)


def startClient(serverIp, serverPort, pipeNum, dataItemNum):
    asyncio.get_event_loop().run_until_complete(
        hello(('ws://%s:%d/' % (serverIp, serverPort)), pipeNum, dataItemNum))
        

if __name__ == "__main__":
    pipeNum = int(sys.argv[1])
    dataItemNum = int(sys.argv[2])
    startClient("192.168.9.72", 1999, pipeNum, dataItemNum)