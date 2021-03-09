// Package httpclient ...
package httpclient

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Get(link string, values map[string]string) ([]byte, error) {
	if !strings.Contains(link, "http://") && !strings.Contains(link, "https://") {
		link = "http://" + link
	}
	baseURL, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	for key, value := range values {
		params.Add(key, value)
	}

	baseURL.RawQuery = params.Encode()

	link = baseURL.String()

	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func Post(link string, values map[string]string) ([]byte, error) {
	if !strings.Contains(link, "http://") && !strings.Contains(link, "https://") {
		link = "http://" + link
	}

	params := url.Values{}
	for key, value := range values {
		params.Add(key, value)
	}

	resp, err := http.PostForm(link, params)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
