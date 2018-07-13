#!/usr/bin/env python3

import asyncio
import websockets
import time
import sys
import array
import struct
import json

async def hello(uri):
    # recvd = 0
    async with websockets.connect(uri) as websocket:

        jsonData = {
            "serverIP": "192.168.9.72",
            "subscribe": [
                {
                    "pipe": 1,
                    "count": 2000,
                    "startFreq": 100,
                    "stopFreq": 2000
                }
                ,
                {
                    "pipe": 2,
                    "count": 2000,
                    "startFreq": 100,
                    "stopFreq": 2000
                },
                {
                    "pipe": 3,
                    "count": 2000,
                    "startFreq": 100,
                    "stopFreq": 2000
                },
                {
                    "pipe": 4,
                    "count": 2000,
                    "startFreq": 100,
                    "stopFreq": 2000
                }
            ]
        }

        jsonStr = json.dumps(jsonData)
        await websocket.send(jsonStr)
        print("json data sent")
        received = 0
        try:
            while True:
                msg = await websocket.recv()
                # jsonmsg = ""
                print("type of msg:", type(msg))
                # strMsg = msg.decode('utf-8')
                # print("strMsg:", strMsg)
                jsonMsg = json.loads(msg)
                # print("client received msg", len(msg))
                print("pipe:", jsonMsg['pipe'])
                print("count:", jsonMsg['count'])
                # print("data:", jsonMsg['data'])
                received += 1
                # da = array.array('d')
                # da.frombytes(msg)
                # da.frombytes(msg)
                # print("received", len(da))
        except:
            print("exception when receiving data after received ", received, "msg(s)")


def startClient(serverIp, serverPort):
    asyncio.get_event_loop().run_until_complete(
        hello(('ws://%s:%d/' % (serverIp, serverPort))))
        

if __name__ == "__main__":
    # pipeNum = int(sys.argv[1])
    # dataItemNum = int(sys.argv[2])
    startClient("127.0.0.1", 1999)