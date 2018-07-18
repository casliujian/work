package main

import "encoding/json"

//将接收到的请求数据解析成map,然后再根据msgType的值做进一步解析
func bytes2map(bytes []byte) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(bytes, &m)
	return m
}

func interface2int(v interface{}) int {
	return int(v.(float64))
}

func interface2string(v interface{}) string {
	return v.(string)
}

