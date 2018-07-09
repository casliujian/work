#!/usr/bin/env python3

import asyncio
import websockets
import time
import sys
import array
import struct

recvno = 40000
channel = int(sys.argv[1])
total = int(sys.argv[2])

async def hello(uri):
    recvd = 0
    async with websockets.connect(uri) as websocket:
        # while True:
        # await websocket.send("Hello world!")
            # time.sleep(1)
        # websockets.recv
        # websocket.send((channel, total))
        channelb = channel.to_bytes(4, byteorder='big',signed=False)
        await websocket.send(channelb)
        totalb = total.to_bytes(4, byteorder='big',signed=False)
        await websocket.send(totalb)
        msg = await websocket.recv()
        # for msg in websocket:
        if recvd != recvno:
            print('received:', struct.unpack('d',msg))
            recvd += 1
        


asyncio.get_event_loop().run_until_complete(
    hello('ws://127.0.0.1:10019'))