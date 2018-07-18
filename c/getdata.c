
// The example of get data
// gcc getdata.c -o getdata -I../include -L../lib -Wl,-R../lib -lspuapi -lcurl

#include <stdio.h>
#include <time.h>
#include "spuapi.h"

int main()
{
    // 创建查询句柄
    SpuHandle h = spuNew();
    
    // 查询频谱数据
    // if(!spuQueryAll(h, time(0)-86400, time(0)))
    // if(!spuQueryAlm(h, 1234)
    if(!spuQuerySam(h, time(0)-86400, time(0)))
    {
        fprintf(stderr, "Execute query error: %s\n", spuStrError(h));
        spuFree(h);
        return 1;
    }
    
    // 获取记录个数
    int count = spuCount(h);
    printf("Count = %d\n", count);
    
    // 读取数据
    char buf[1024*1024];
    int cnt = 0;
    while(spuNext(h))
    {
        int ret = spuFetch(h, buf, sizeof(buf));
        if(ret < 0)
        {
            fprintf(stderr, "Fetch data error: %s\n", spuStrError(h));
            break;
        }
        printf("cnt%d, datlen=%d\n", ++cnt, ret);
    }
    
    // 释放查询句柄
    spuFree(h);
    
    return 0;
}
