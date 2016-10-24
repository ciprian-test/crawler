package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const defaultWorkersCount = 2
const defaultMaxLevel = 5

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) < 2 {
		fmt.Println("The start URL is missing")
		os.Exit(1)
	}

	urlArg := os.Args[1]
	if !strings.HasPrefix(urlArg, "http") {
		// Allow simple start urls and default using http
		urlArg = "http://" + urlArg
	}

	startURL, err := validateStartURL(urlArg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scraper := NewScraper(startURL, getWorkersCount(), getMaxLevelCount())
	scraper.Run()
	scraper.PrintJSON()
}

// Basic validation of the start URL
func validateStartURL(urlArg string) (*url.URL, error) {
	if len(urlArg) > 2000 || len(urlArg) < 11 {
		return nil, errors.New("Invalid start URL")
	}

	startURL, err := url.Parse(urlArg)
	if err != nil {
		return nil, fmt.Errorf("Invalid start URL: %s", err)
	}

	if !startURL.IsAbs() {
		return nil, errors.New("The start URL is relative")
	}

	return startURL, nil
}

func getWorkersCount() int {
	workers := defaultWorkersCount
	if len(os.Args) >= 3 {
		workersNumber, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Println("Invalid workers parameter")
			os.Exit(1)
		}
		workers = int(workersNumber)
	}
	return workers
}

func getMaxLevelCount() int {
	maxLevel := defaultMaxLevel
	if len(os.Args) >= 4 {
		maxLevelNumber, err := strconv.ParseInt(os.Args[3], 10, 64)
		if err != nil {
			fmt.Println("Invalid max level parameter")
			os.Exit(1)
		}
		maxLevel = int(maxLevelNumber)
	}
	return maxLevel
}
