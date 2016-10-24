package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
)

// Serve a few static pages located in the 'test_data' folder.
func TestMain(m *testing.M) {
	server := http.Server{
		Addr:    ":9000",
		Handler: http.FileServer(http.Dir("test_data")),
	}

	go func() {
		server.ListenAndServe()
	}()

	os.Exit(m.Run())
}

// Scrape the localhost domain, starting with 'index.html' and only 2 levels.
// Match the results with the 'result.json' file.
func TestIntegration(t *testing.T) {
	startURL, _ := url.Parse("http://localhost:9000/index.html")
	scraper := NewScraper(startURL, 2, 2)
	scraper.Run()

	results, err := scraper.GetJSON()
	if err != nil {
		t.Fatal("Error printing JSON")
	}

	resultData, err := ioutil.ReadFile("test_data/result.json")
	if err != nil {
		t.Fatal("Error reading expected data from file")
	}

	if string(resultData) != string(results) {
		t.Logf("Incorrect output: expected = \n%s\n\nactual = \n%s", resultData, results)
		t.Fail()
	}
}
