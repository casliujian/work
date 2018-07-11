#!/usr/bin/env python3

import asyncio
import websockets
import time
import sys
import array
import struct
import json

recvno = 40000
channel = int(sys.argv[1])
total = int(sys.argv[2])

async def hello(uri):
    recvd = 0
    async with websockets.connect(uri) as websocket:

        # channelb = channel.to_bytes(4, byteorder='little',signed=False)
        # await websocket.send(channelb)
        # totalb = total.to_bytes(4, byteorder='little',signed=False)
        # await websocket.send(totalb)

        jsonData = json.dumps({'pipeNum': 1, 'dataItemNum': 2000})
        # await websocket.send("helloworld")
        await websocket.send(jsonData)
        print("json data sent")

        msg = await websocket.recv()
        # for msg in websocket:
        if recvd != recvno:
            print('received:', msg)
            recvd += 1
        


asyncio.get_event_loop().run_until_complete(
    hello('ws://127.0.0.1:1999/'))