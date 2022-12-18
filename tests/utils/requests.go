package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const (
	baseUrl = "http://localhost:8080/"
)

func encodeStruct(value any) *bytes.Buffer {
	jsonValue, _ := json.Marshal(value)
	return bytes.NewBuffer(jsonValue)
}

func decodeBody(resBody io.ReadCloser, target any) error {
	return json.NewDecoder(resBody).Decode(&target)
}

func Get(path string, response any) error {
	res, err := http.Get(baseUrl + path)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return decodeBody(res.Body, &response)
}

func Post(path string, value any, response any) error {
	res, err := http.Post(baseUrl+path, "application/json", encodeStruct(value))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return decodeBody(res.Body, &response)
}
