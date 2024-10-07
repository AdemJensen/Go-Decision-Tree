package utils

import "encoding/json"

func JsonPretty(data interface{}) string {
	res, _ := json.MarshalIndent(data, "", "  ")
	return string(res)
}

func Json(data interface{}) string {
	res, _ := json.Marshal(data)
	return string(res)
}
