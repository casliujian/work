package main

import (
    "fmt"
    "encoding/json"
)

type Data struct {
    Votes *Votes `json:"votes"`
    Count string `json:"count,omitempty"`
}

type Votes struct {
    OptionA string `json:"option_A"`
}

type TestJson struct {
    Num int `json:"num"`
    Str []string `json:"str"`
}

func main() {
    //s := `{ "votes": { "option_A": "3" } }`
    test := `{"num": 123,"str": ["helloworld", "yes"]}`
    //data := &Data{
    //    Votes: &Votes{},
    //}
    testJson := new(TestJson)
    json.Unmarshal([]byte(test), testJson)
    fmt.Printf("num: %d, str: %s, %s\n", testJson.Num, testJson.Str[0], testJson.Str[1])
	//err := json.Unmarshal([]byte(s), data)
	//fmt.Println(err)
	//fmt.Println(data.Votes)
	//s2, _ := json.Marshal(data)
	//fmt.Println(string(s2))
	//data.Count = "2"
	//s3, _ := json.Marshal(data)
	//fmt.Println(string(s3))
	//println(data.Count)
}
