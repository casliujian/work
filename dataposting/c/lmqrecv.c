
// The example of subscribe message

#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include "lmq.h"

void main()
{
    // Create Lmq socket
    int sock = lmqSocket(LMQ_SUB);
    if(sock < 0)
    {
        fprintf(stderr, "Create socket error: %s\n", lmqStrError(lmqErrno()));
        return;
    }
    
    // Subscribe message
    if(lmqSub(sock, 1000) < 0)
    {
        fprintf(stderr, "Subscribe message error: %s\n", lmqStrError(lmqErrno()));
        return;
    }
    if(lmqSub(sock, 1001) < 0)
    {
        fprintf(stderr, "Subscribe message error: %s\n", lmqStrError(lmqErrno()));
        return;
    }
    
    // Receive message
    int cnt = 0;
    while(1)
    {
        unsigned int msgID;
        char buf[1024];
        bzero(buf, sizeof(buf));
        int ret = lmqRecv(sock, &msgID, buf, sizeof(buf), 0);
        if(ret < 0)
        {
            fprintf(stderr, "Receive message error: %s\n", lmqStrError(lmqErrno()));
            return;
        }
        
        printf("cnt = %d id = %d len = %d msg = %s\n", 
                ++cnt, msgID, ret, buf);
    }
}
