package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	authUrl      = "http://localhost:8110/"
	firestoreUrl = "http://localhost:8080/"
)

func encodeStruct(value any) io.Reader {
	if value == nil {
		return nil
	}
	jsonValue, _ := json.Marshal(value)
	return bytes.NewBuffer(jsonValue)
}

func decodeBody(resBody io.ReadCloser, target any) error {
	return json.NewDecoder(resBody).Decode(&target)
}

func Get(path string, response any) error {
	res, err := http.Get(firestoreUrl + path)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return decodeBody(res.Body, &response)
}

type UrlParameter struct {
	Key   string
	Value string
}

type Url struct {
	Path   string
	Params []UrlParameter
}

func GetWithToken(url Url, loginDetails LoginUser, response any) error {
	return requestWithToken(http.MethodGet, url, nil, loginDetails, response)
}

func PostWithToken(url Url, data any, loginDetails LoginUser, response any) error {
	return requestWithToken(http.MethodPost, url, data, loginDetails, response)
}

func requestWithToken(method string, url Url, data any, loginDetails LoginUser, response any) error {
	token, err := GenerateIdToken(loginDetails)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, firestoreUrl+url.Path, encodeStruct(data))
	if err != nil {
		fmt.Printf("Creating %v request failed. Given url: %v\n", method, url.Path)
		return err
	}

	if len(url.Params) > 0 {
		q := req.URL.Query()
		for _, param := range url.Params {
			q.Add(param.Key, param.Value)
		}
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Set(os.Getenv("AUTHHEADER"), token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("%v request failed. Given url: %v. %v", method, url, err)
		return err
	}
	defer res.Body.Close()

	return decodeBody(res.Body, &response)
}

func Post(path string, value any, response any) error {
	res, err := http.Post(firestoreUrl+path, "application/json", encodeStruct(value))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return decodeBody(res.Body, &response)
}

func PostAuth(path string, value any, response any) error {
	res, err := http.Post(authUrl+path, "application/json", encodeStruct(value))
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
