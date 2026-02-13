package utils

import "encoding/json"

func DecodeJSON(data string, result any) error {
	return json.Unmarshal([]byte(data), result)
}
