package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func Delete(path string) error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		fmt.Println("Creating delete request failed. Given url:", path)
		return err
	}

	_, err = client.Do(req)
	if err != nil {
		fmt.Println("Delete failed. Given url:", path)
		return err
	}

	return nil
}
