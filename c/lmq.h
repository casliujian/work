
/**
 @file lmq.h
 @brief Light message queue
 @author lxc
 @date 2018/4/3 21:01:30
 @version 1.0
*/

#ifndef _LMQ_H__
#define _LMQ_H__

#define LMQ_SUB_ENDPOINT        "ipc:///tmp/lmq.ipc.sub"    ///< Endpoint of subscribe
#define LMQ_PUB_ENDPOINT        "ipc:///tmp/lmq.ipc.pub"    ///< Endpoint of publish
#define LMQ_PUB_ENDPOINT_TCP    "tcp://*:2727"              ///< Endpoint of publish for tcp
#define LMQ_RCVBUF              (1*1024*1024)               ///< Size of the receive buffer, in bytes
#define LMQ_SNDBUF              (1*1024*1024)               ///< Size of the send buffer, in bytes
#define LMQ_MAX_MSGSZ           (1024*1024)                 ///< Maximum message size, in bytes
#define LMQ_PUB                 32                          ///< Scoket for publish message
#define LMQ_SUB                 33                          ///< Scoket for subscribe message
#define LMQ_DONTWAIT            1                           ///< Specifies that the operation should be performed in non-blocking mode

/// Lmq message struct
typedef struct LmqMsg
{
    unsigned int id;                                        ///< Message ID
    unsigned char buf[LMQ_MAX_MSGSZ];                       ///< Message data
}LmqMsg;

#ifdef __cplusplus
extern "C" {
#endif

/**
 @brief Create Lmq socket
 @param[in] type Socket type, LMQ_SUB or LMQ_PUB
 @return If the function succeeds file descriptor of the new socket is returned. 
         Otherwise, -1 is returned and errno is set.
 */
int lmqSocket(int type);

/**
 @brief Create Lmq socket(SUB) for tcp
 @param[in] broker Broker address
 @return If the function succeeds file descriptor of the new socket is returned. 
         Otherwise, -1 is returned and errno is set.
 */
int lmqSocketTcp(const char* broker);

/**
 @brief Close Lmq socket
 @param[in] sock Lmq socket
 @return If the function succeeds zero is returned. 
         Otherwise, -1 is returned and errno is set.
 */
int lmqClose(int sock);

/**
 @brief Subscribe message
 @param[in] sock Lmq socket
 @param[in] msgID Message ID
 @return If the function succeeds zero is returned. 
         Otherwise, -1 is returned and errno is set.
 */
int lmqSub(int sock, unsigned int msgID);

/**
 @brief Unsubscribe message
 @param[in] sock Lmq socket
 @param[in] msgID Message ID
 @return If the function succeeds zero is returned. 
         Otherwise, -1 is returned and errno is set.
 */
int lmqUnsub(int sock, unsigned int msgID);

/**
 @brief Send message
 @param[in] sock Lmq socket
 @param[in] msgID MessageID
 @param[in] buf Message buffer
 @param[in] len Message length
 @param[in] flag The operation flag, 0 or LMQ_DONTWAIT
 @return If the function succeeds, the number of bytes in the message is returned. 
         Otherwise, -1 is returned and errno is set.
 */
int lmqSend(int sock, unsigned int msgID, const void* buf, int len, int flag);

/**
 @brief Receive message
 @param[in] sock Lmq socket
 @param[out] msgID MessageID
 @param[out] buf Receive buffer
 @param[in] len The length of receive buffer
 @param[in] flag The operation flag, 0 or LMQ_DONTWAIT
 @return If the function succeeds number of bytes in the message is returned. 
         Otherwise, -1 is returned and errno is set
 */
int lmqRecv(int sock, unsigned int* msgID, void* buf, int len, int flag);

/**
 @brief Retrieve value of errno for the calling thread
 @return The value of the errno.
 */
int lmqErrno();

/**
 @brief Convert an error number into human-readable string
 @return Return error message string.
 */
const char* lmqStrError(int err);

#ifdef __cplusplus
}
#endif

#endif
