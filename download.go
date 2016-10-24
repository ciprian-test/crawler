package main

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// Make a http request and return the response body.
func getPageContent(pageURL string) (string, error) {
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", errors.New("Invalid response status code")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("Error reading body")
	}

	return string(body), nil
}
