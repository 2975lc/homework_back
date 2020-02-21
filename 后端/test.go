package main

import (
	"encoding/json"
	"fmt"
)
type IT struct {
	Company  string
	Subjects string
	IsOk     bool
	Price    float64
}
func main (){
	t1 := IT{"tencent", "develop", false, 12000.0}
	b, err := json.Marshal(t1)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Println(string(b))
}