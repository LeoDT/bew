package bew

import (
	"encoding/json"
	"net/http"
)

type BaseDict interface {
	Get(string) string
	Add(string, string)
	Del(string)
	Set(string, string)
}

type JsonDict map[string]interface{}

type FileDict BaseDict

func (j *JsonDict) Parse(r *http.Request) error {
	body := make([]byte, r.ContentLength)
	_, err := r.Body.Read(body)

	if err != nil {
		return nil
	}
	
	var json_map map[string]interface{}

	err = json.Unmarshal(body, &json_map)

	if err != nil {
		return err
	}

	*j = json_map

	return nil
}












