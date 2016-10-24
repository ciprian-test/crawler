package main

import (
	"net/url"
	"regexp"
	"strings"
)

const regexpScript string = `<script [^>]*?src=[ ]*['"]?([^>'"]+)['"]?[^>]*>`
const regexpImage string = `<img [^>]*?src=[ ]*['"]?([^>'"]+)['"]?[^>]*>`
const regexpCSSLink string = `<link [^>]*href=[ ]*['"]?([^>'"]+)['"]?[^>]*>`
const regexpLink string = `<a [^>]*href=[ ]*['"]?([^>'"]+)['"]?[^>]*>`

var rxScript = regexp.MustCompile(regexpScript)
var rxImage = regexp.MustCompile(regexpImage)
var rxCSSLink = regexp.MustCompile(regexpCSSLink)
var rxLink = regexp.MustCompile(regexpLink)

// Extract information about different types of links from a page body.
// Identifies scripts, images, CSS links and internal and external links.
func extractDataFromBody(host string, pageURL string, body string) *Page {
	parsedURL, _ := url.Parse(pageURL)
	oneLineBody := strings.Replace(body, "\n", " ", -1)

	scripts := extractGenericURL(parsedURL, oneLineBody, rxScript)
	images := extractGenericURL(parsedURL, oneLineBody, rxImage)
	cssLinks := extractGenericURL(parsedURL, oneLineBody, rxCSSLink)
	links := extractGenericURL(parsedURL, oneLineBody, rxLink)

	internalLinks, externalLinks := identifyLinks(host, links)

	return &Page{
		URL:           pageURL,
		Scripts:       scripts,
		Images:        images,
		CSSLinks:      cssLinks,
		InternalLinks: internalLinks,
		ExternalLinks: externalLinks,
		Scraped:       true,
	}
}

// Apply a Regexp on the page body and return the matched URLs.
func extractGenericURL(pageURL *url.URL, body string, rx *regexp.Regexp) []string {
	urls := []string{}

	submatches := rx.FindAllStringSubmatch(body, -1)
	if submatches != nil {
		for _, submatch := range submatches {
			if len(submatch) >= 1 {
				subURL, err := absoluteURL(pageURL, strings.TrimSpace(submatch[1]))
				if err == nil {
					urls = append(urls, subURL)
				}
			}
		}
	}

	return urls
}

// Resolve links found in the page relative to the parent page URL.
func absoluteURL(pageURL *url.URL, childURL string) (string, error) {
	subURL, err := url.Parse(childURL)
	if err != nil {
		return "", err
	}

	return pageURL.ResolveReference(subURL).String(), nil
}

// Split links into 'internal' and 'external'.
func identifyLinks(host string, links []string) ([]string, []string) {
	internalLinks := []string{}
	externalLinks := []string{}

	for _, link := range links {
		matches, err := regexp.MatchString("https?://"+host, link)
		if err != nil {
			continue
		}

		if matches {
			internalLinks = append(internalLinks, link)
		} else {
			externalLinks = append(externalLinks, link)
		}
	}

	return internalLinks, externalLinks
}
