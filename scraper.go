package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"
)

// Scraper structure for storing scraper data
type Scraper struct {
	sync.RWMutex
	host          string
	Pages         map[string]*Page `json:"pages"`
	remainingURLs map[string]*URLInfo
	pendingURLs   map[string]*URLInfo
	workers       int
	maxLevel      int
}

// Page structure for storing URLs identified in a page body
type Page struct {
	URL     string `json:"url"`
	Scraped bool   `json:"scraped"`

	Scripts       []string `json:"scripts"`
	Images        []string `json:"images"`
	CSSLinks      []string `json:"cssLinks"`
	InternalLinks []string `json:"internalLinks"`
	ExternalLinks []string `json:"externalLinks"`
}

// URLInfo url info
type URLInfo struct {
	URL   string
	Level int
}

// NewScraper creates a new scraper.
func NewScraper(startURL *url.URL, workers int, maxLevel int) *Scraper {
	return &Scraper{
		host:  startURL.Host,
		Pages: map[string]*Page{},
		remainingURLs: map[string]*URLInfo{
			startURL.String(): &URLInfo{
				URL:   startURL.String(),
				Level: 1,
			},
		},
		pendingURLs: map[string]*URLInfo{},
		workers:     workers,
		maxLevel:    maxLevel,
	}
}

// Run starts the scraping.
func (s *Scraper) Run() {
	var wg sync.WaitGroup

	for i := 1; i <= s.workers; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			fmt.Printf("Started worker %d\n", i)
			for {
				stop := s.scrape()
				if stop {
					break
				}
			}
			fmt.Printf("Stopped worker %d\n", i)
		}(i)
	}

	wg.Wait()
}

func (s *Scraper) scrape() bool {
	s.RLock()
	if len(s.remainingURLs) == 0 && len(s.pendingURLs) == 0 {
		s.RUnlock()
		return true
	}
	s.RUnlock()

	s.Lock()

	var urlInfo *URLInfo
	for _, info := range s.remainingURLs {
		urlInfo = info
		break
	}

	if urlInfo != nil {
		delete(s.remainingURLs, urlInfo.URL)
		s.pendingURLs[urlInfo.URL] = urlInfo
		s.Unlock()

		s.processURL(urlInfo)
	} else {
		s.Unlock()
		time.Sleep(200 * time.Millisecond)
	}

	return false
}

func (s *Scraper) processURL(urlInfo *URLInfo) {
	fmt.Printf("Process: %s\n", urlInfo.URL)

	body, err := getPageContent(urlInfo.URL)
	if err != nil {
		s.Lock()
		defer s.Unlock()

		page := &Page{
			URL:     urlInfo.URL,
			Scraped: false,
		}
		s.Pages[urlInfo.URL] = page
		delete(s.pendingURLs, urlInfo.URL)
		return
	}

	page := extractDataFromBody(s.host, urlInfo.URL, body)

	s.Lock()
	defer s.Unlock()

	s.Pages[urlInfo.URL] = page
	delete(s.pendingURLs, urlInfo.URL)

	if urlInfo.Level < s.maxLevel {
		s.addAdditionalURLs(urlInfo, page)
	}
}

func (s *Scraper) addAdditionalURLs(parentURL *URLInfo, page *Page) {
	for _, internalLink := range page.InternalLinks {
		if _, ok := s.Pages[internalLink]; ok {
			continue
		}
		if _, ok := s.remainingURLs[internalLink]; ok {
			continue
		}
		if _, ok := s.pendingURLs[internalLink]; ok {
			continue
		}
		s.remainingURLs[internalLink] = &URLInfo{
			URL:   internalLink,
			Level: parentURL.Level + 1,
		}
	}
}

// Print outputs data to console
func (s *Scraper) Print() {
	fmt.Println("\n\nResults")

	for _, page := range s.Pages {
		if !page.Scraped {
			continue
		}

		fmt.Printf("URL: %s\n\n", page.URL)

		printList(page.Scripts, "Scripts")
		printList(page.CSSLinks, "CSS Links")
		printList(page.Images, "Images")
		printList(page.ExternalLinks, "External Links")
		printList(page.InternalLinks, "Internal Links")
	}
}

func printList(list []string, title string) {
	fmt.Println("  " + title + ":")
	for _, item := range list {
		fmt.Printf("    %s\n", item)
	}
	fmt.Println("")
}

// GetJSON returns a JSON with the results
func (s *Scraper) GetJSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

// PrintJSON outputs to the console a JSON with the results
func (s *Scraper) PrintJSON() {
	result, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(result))
}
