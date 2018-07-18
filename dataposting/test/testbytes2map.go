package main

import (

	"encoding/json"
)

func main() {
	jsonstr := []byte("{\"item1\": 1, \"item2\": \"hello\"}")

	var m map[string]interface{}

	//println(reflect.TypeOf(jsonstr))
	json.Unmarshal(jsonstr, &m)

	//fmt.Printf("\n",int(m["item1"].(float64)))
	println(int(m["item1"].(float64)))
	println(m["item2"].(string))



}
