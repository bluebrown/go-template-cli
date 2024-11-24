package textfunc

import "encoding/json"

func FromJson(s string) any {
	var v any
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		return nil
	}
	return v
}
