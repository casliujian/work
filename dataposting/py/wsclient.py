#!/usr/bin/env python3

import asyncio
import websockets
import time
import sys
import array
import struct
import json

async def hello(uri, cno):
    # recvd = 0
    async with websockets.connect(uri) as websocket:

        jsonData = {
            "msgType": 1,
            "serverIP": "192.168.31.200",
            "subscribe": [
                {
                    "pipe": 0,
                    "count": 2000,
                    "startFreq": 100,
                    "stopFreq": 2000
                }
                # ,
                # {
                #     "pipe": 1,
                #     "count": 2000,
                #     "startFreq": 100,
                #     "stopFreq": 2000
                # },
                # {
                #     "pipe": 2,
                #     "count": 2000,
                #     "startFreq": 100,
                #     "stopFreq": 2000
                # },
                # {
                #     "pipe": 3,
                #     "count": 2000,
                #     "startFreq": 100,
                #     "stopFreq": 2000
                # }
            ]
        }

        jsonStr = json.dumps(jsonData)
        await websocket.send(jsonStr)
        # print("json data sent")
        received = {}
        receivedNo = 0
        try:
            while True:
                msg = await websocket.recv()
                # jsonmsg = ""
                # print("type of msg:", type(msg))
                print("size of msg:", len(msg))
                # print("msg:\n", msg)
                jsonData = json.loads(msg.decode('utf-8'))
                receivedNo += 1
                # received[jsonData["pipe"]] += 1
                # pipeNum = jsonData['pipe']
                # if pipeNum not in received:
                #    received[pipeNum] = 1
                # else:
                #    received[pipeNum] += 1
        except Exception:
            print("client", cno, "received ", receivedNo, "msg(s)")
            # print('exception:', e)
            # print("client", cno, "received for each pipe", received)

def startClient(serverIp, serverPort, cno):
    asyncio.get_event_loop().run_until_complete(
        hello(('ws://%s:%d/' % (serverIp, serverPort)), cno))
        

if __name__ == "__main__":
    # pipeNum = int(sys.argv[1])
    # dataItemNum = int(sys.argv[2])
    startClient("127.0.0.1", 1999, 0)