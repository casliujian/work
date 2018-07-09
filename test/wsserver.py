#!/usr/bin/env python3

import asyncio
import websockets
import time
import socket
import queue
import threading
import array
import struct

registry = {}
class Register:
    def __init__(self, interval, total):
        # self.register_count = 1
        self.interval = interval
        self.total = total
        self.data_queue = queue.Queue(total)
    # def incrCount(self):
    #     self.register_count += 1
    # def decrCount(self):
    #     self.register_count -= 1




class ReceiveDataThread(threading.Thread):
    def __init__(self, resisty):
        threading.Thread.__init__(self)
        datasocket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        datasocket.connect(('127.0.0.1', 2000))
        self.sockt = datasocket
        self.registry = resisty
    def run(self):
        channel_byte_len = 4*4
        data_byte_len = 65536*8
        # total_byte_len = channel_byte_len + data_byte_len
        # received_byte_len = 0
        # current_channel = None
        round = 0
        start = time.time()
        finish = time.time()
        try:
            while (finish - start) < 4:
                # Parse the channel index
                received_channel_byte_len = 0
                channel_bytes = bytearray(channel_byte_len)
                while received_channel_byte_len < channel_byte_len:
                    received = self.sockt.recv(channel_byte_len - received_channel_byte_len)
                    received_len = len(received)
                    channel_bytes[received_channel_byte_len : received_channel_byte_len + received_len] = received
                    received_channel_byte_len += received_len
                # Calculate the channel index from bytes, big endian, unsigned
                channel_idx = int.from_bytes(channel_bytes[0 : channel_byte_len], byteorder='big', signed=False)

                # Now parse the float data
                received_data_byte_len = 0
                data_bytes = bytearray(data_byte_len)
                while received_data_byte_len < data_byte_len:
                    received = self.sockt.recv(data_byte_len - received_data_byte_len)
                    received_len = len(received)
                    data_bytes[received_data_byte_len : received_len + received_len] = received
                    received_data_byte_len += received_len
                # Conver byte array into double array 
                data_array = array.array('d')
                data_array.frombytes(data_bytes[0:data_byte_len])
                # for d in data_array:
                #     print(d)
                round += 1
                if channel_idx in self.registry:
                    for r in self.registry[channel_idx]:
                        current_data_idx = 0
                        for i in range(r.total):
                            q = r.data_queue
                            q.put(1+data_array[current_data_idx])
                            current_data_idx += r.interval
                finish = time.time()
        except KeyboardInterrupt:
            finish = time.time()
            print('rounds per second:', round/((finish - start)))
            self.sockt.close()
        print('rounds per second:', round/((finish - start)))
        self.sockt.close()
            # received_bytes = self.sockt.recv(total_byte_len)
            # received_byte_len += (len(received_bytes)+received_byte_len)%total_byte_len






def register(request):
    channel_idx, total = request
    # q = queue.Queue(65536)
    interval = 65536//total
    r = Register(interval, total)
    if channel_idx not in registry:
        
        registry[channel_idx] = [r]
        return r.data_queue
    else:
        # for r in registry[channel_idx]:
        #     if dot == r.interested_dot:
        #         r.incrCount()
        #         return r
        registry[channel_idx].append(r)
        return r.data_queue


# async def wssend(websocket, data):
#     pass

async def echo(ws, path):
    channelb = await ws.recv()
    channel = int.from_bytes(channelb, byteorder='big', signed=False)
    totalb = await ws.recv()
    total = int.from_bytes(totalb, byteorder='big',signed=False)
    print('A new client want channel', channel, 'for', total, 'items')
    q = register((channel,total))
    async while True:    
        data_item = await q.get(block=True)
        print('websocket server sending', data_item)
        # data_item = 1.0001
        await ws.send(struct.pack('d',data_item))



received_thread = ReceiveDataThread(registry)
received_thread.start()
# received_thread.join()

asyncio.get_event_loop().run_until_complete(
    websockets.serve(echo, '127.0.0.1', 10019))
try:
    asyncio.get_event_loop().run_forever()
except KeyboardInterrupt:
    pass
    # datasocket.close()