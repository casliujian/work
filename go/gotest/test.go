package main

var testmap = map[int]TestStruct{}
type TestStruct struct {
	Item1 int
	Item2 string
}

func main() {
	testmap[1] = TestStruct{1, "2"}
	ts := (testmap[1])
	ts.Item1 = 3
	println(testmap[1].Item1)
	println(testmap[1].Item2)
}
