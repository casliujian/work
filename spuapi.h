
/**
 @file spuapi.h
 @brief 频谱数据查询接口库
 @author lxc
 @date 2018/5/12 16:38:58
 @version 1.0
*/

#ifndef _SPUAPI_H__
#define _SPUAPI_H__

typedef void* SpuHandle;        ///< 查询句柄


#ifdef __cplusplus
extern "C" {
#endif

/**
 @brief 创建查询句柄
 @return 成功时返回查询句柄，失败时返回空指针。
 */
SpuHandle spuNew();

/**
 @brief 释放查询句柄
 @param[in] handle 查询句柄
 @return 无
 */
void spuFree(SpuHandle handle);

/**
 @brief 查询全态频谱数据
 @param[in] handle 查询句柄
 @param[in] btime 开始时间，基于1970-01-01的累积秒。
 @param[in] etime 结束时间，基于1970-01-01的累积秒。
 @return 成功时返回1，失败时返回0。
 */
int spuQueryAll(SpuHandle handle, long btime, long etime);

/**
 @brief 查询抽样频谱数据
 @param[in] handle 查询句柄
 @param[in] btime 开始时间，基于1970-01-01的累积秒。
 @param[in] etime 结束时间，基于1970-01-01的累积秒。
 @return 成功时返回1，失败时返回0。
 */
int spuQuerySam(SpuHandle handle, long btime, long etime);

/**
 @brief 查询告警频谱数据
 @param[in] handle 查询句柄
 @param[in] almTime 告警时间。
 @return 成功时返回1，失败时返回0。
 */
int spuQueryAlm(SpuHandle handle, long almTime);

/**
 @brief 获取记录个数
 @param[in] handle 查询句柄
 @return 返回记录个数
 */
int spuCount(SpuHandle handle);

/**
 @brief 设置为第一条记录
 @param[in] handle 查询句柄
 @return 成功时返回1，失败时返回0。
 */
int spuFirst(SpuHandle handle);

/**
 @brief 设置为最后一条记录
 @param[in] handle 查询句柄
 @return 成功时返回1，失败时返回0。
 */
int spuLast(SpuHandle handle);

/**
 @brief 设置为后一条记录
 @param[in] handle 查询句柄
 @return 成功时返回1，失败时返回0。
 */
int spuNext(SpuHandle handle);

/**
 @brief 设置为前一条记录
 @param[in] handle 查询句柄
 @return 成功时返回1，失败时返回0。
 */
int spuPrev(SpuHandle handle);

/**
 @brief 设置为指定记录
 @param[in] handle 查询句柄
 @param[in] index 记录索引，从0开始。
 @return 成功时返回1，失败时返回0。
 */
int spuSeek(SpuHandle handle, int index);

/**
 @brief 读取数据
 @param[in] handle 查询句柄
 @param[out] buf 数据缓存
 @param[in] len 缓存长度
 @return 读取成功返回数据长度，失败时返回-1。
 */
int spuFetch(SpuHandle handle, void* buf, int len);

/**
 @brief 获取错误说明
 @param[in] handle 查询句柄
 @param[out] buf 数据缓存
 @param[in] len 缓存长度
 @return 错误说明
 */
const char* spuStrError(SpuHandle handle);

#ifdef __cplusplus
}
#endif

#endif
