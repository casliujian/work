from nanomsg import Socket, PUB
import sys
import array


def server(url):
    sock = Socket(PUB)
    sock.bind(url)

    
    channel1 = b'\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
    channel2 = b'\x02\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
    channel3 = b'\x03\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
    channel4 = b'\x04\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'

    doublearray = array.array('d')
    for i in range(65536):
        doublearray.insert(i,0.000001+i)

    doublebytes = doublearray.tobytes()

    # databytes = doublearray.tobytes()
    # pipes1 = channel1.decode('utf-8')
    # pipes2 = channel2.decode('utf-8')
    # pipes3 = channel3.decode('utf-8')
    # pipes4 = channel4.decode('utf-8')
    # datas = databytes.decode('utf-8')
    # pipedatas1 = pipes1+datas
    # pepedata1 = pipedatas1.encode('utf-8')
    # pipedatas2 = pipes2+datas
    # pepedata2 = pipedatas2.encode('utf-8')
    # pipedatas3 = pipes3+datas
    # pepedata3 = pipedatas3.encode('utf-8')
    # pipedatas4 = pipes4+datas
    # pepedata4 = pipedatas4.encode('utf-8')
    while True:
        sock.send(b''.join([channel1,doublebytes]))
        sock.send(b''.join([channel2,doublebytes]))
        sock.send(b''.join([channel3,doublebytes]))
        sock.send(b''.join([channel4,doublebytes]))


if __name__ == "__main__":
    server(sys.argv[1])