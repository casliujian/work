
// The example of publish message

#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include "lmq.h"

void main()
{
    // Create socket
    int sock = lmqSocket(LMQ_PUB);
    if(sock < 0)
    {
        fprintf(stderr, "Create socket error: %s\n", lmqStrError(lmqErrno()));
        return;
    }
    
    // Send message
    int i = 0;
    for(i = 0; i < 1000; i++)
    {
        char buf[1024];
        sprintf(buf, "hello lmq %d", i+1);
        
        int ret = lmqSend(sock, 1000, buf, strlen(buf), 0);
        if(ret < 0)
        {
            fprintf(stderr, "Send message error: %s\n", lmqStrError(lmqErrno()));
            break;
        }
        
        printf("cnt = %d len = %d msg = %s\n", i+1, ret, buf);
        
        sleep(1);
    }
    
    lmqClose(sock);
}
