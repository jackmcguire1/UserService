package utils

import "encoding/json"

func ToJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
