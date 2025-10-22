package lcutils

import (
	"encoding/json"
	"fmt"
)

func LogJson(v ...any) {
	for _, a := range v {
		bytes, err := json.Marshal(a)
		if err != nil {
			fmt.Println("json.Marshal err:", err)
			return
		}
		fmt.Println(string(bytes))
	}
}
