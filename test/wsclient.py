#!/usr/bin/env python

import asyncio
import websockets
import time

async def hello(uri):
    async with websockets.connect(uri) as websocket:
        while True:
            await websocket.send("Hello world!")
            time.sleep(1)


asyncio.get_event_loop().run_until_complete(
    hello('ws://127.0.0.1:10019'))