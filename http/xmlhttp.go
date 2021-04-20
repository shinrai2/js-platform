package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type XHR struct {
	url string
	header map[string]string
	payload *strings.Reader
	cookies []*http.Cookie
}

func (xhr XHR)SetUrl(u string) {
	xhr.url = u
}

func (xhr XHR)Http(method string) string {
	var (
		err error
		req *http.Request
		res *http.Response
		body []byte
	)
	method = strings.ToUpper(method)
	if strings.Compare(method, "GET") != 0 || strings.Compare(method, "HTTP") != 0 {
		method = "GET"
	}
	client := &http.Client{}
	req, err = http.NewRequest(method, xhr.url, xhr.payload)
	if err != nil {
		fmt.Println(err)
	}
	for _, c := range xhr.cookies {
		req.AddCookie(c)
	}
	for k, v := range xhr.header {
		req.Header.Add(k, v)
	}
	res, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	xhr.cookies = res.Cookies()
	return string(body)
}