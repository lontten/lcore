package utils

import (
	"encoding/json"
	"fmt"
)

func LogJson(v any) {
	bytes, err := json.Marshal(v)
	if err != nil {
		fmt.Println("json.Marshal err:", err)
		return
	}
	fmt.Println(string(bytes))
}
