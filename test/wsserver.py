#!/usr/bin/env python

import asyncio
import websockets

async def echo(websocket, path):
    async for message in websocket:
        print('server received:', message)
        await websocket.send(message)


asyncio.get_event_loop().run_until_complete(
    websockets.serve(echo, '127.0.0.1', 10019))
asyncio.get_event_loop().run_forever()