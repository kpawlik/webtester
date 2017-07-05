package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type param struct {
	Name, Value string
}
type confURL struct {
	URL     string
	Method  string
	Headers []*param
	Params  []*param
}

type data struct {
	Data    []*confURL
	Headers []*param
}

func (r *data) requests() (requests []*http.Request, err error) {
	var (
		req *http.Request
	)
	for _, url := range r.Data {
		req, err = newReq(url, r.Headers)
		if err != nil {
			err = fmt.Errorf("Cannot create request. %v\n", err)
			return
		}
		requests = append(requests, req)
	}
	return
}

func readConfig(filename string) (conf *data, err error) {
	var (
		bytes []byte
	)
	if bytes, err = ioutil.ReadFile(filename); err != nil {
		err = fmt.Errorf("Cannot read config file. %v\n", err)
		return
	}
	if err = json.Unmarshal(bytes, &conf); err != nil {
		err = fmt.Errorf("Cannot decode config file. %v\n", err)
		return
	}
	return
}

func newReq(cfgURL *confURL, headers []*param) (req *http.Request, err error) {
	method := "GET"
	if cfgURL.Method != "" {
		method = cfgURL.Method
	}
	data := url.Values{}
	for _, param := range cfgURL.Params {
		data.Add(param.Name, param.Value)
	}

	req, err = http.NewRequest(method, cfgURL.URL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return
	}
	for _, param := range cfgURL.Headers {
		req.Header.Add(param.Name, param.Value)
	}
	for _, param := range headers {
		req.Header.Add(param.Name, param.Value)
	}

	return
}
