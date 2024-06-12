package tools

import "encoding/json"

func Copy[T any](srv T) (t T) {
	var a T
	b, err := json.Marshal(srv)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &t)
	if err != nil {
		return a
	}
	return
}
