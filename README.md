# Web scraper

## Build and run

Solution implemented using golang (https://golang.org/doc/install).

Build: `make build`. Run `make test` to run the included integration test.
Run: ./build/scraper "<start-url>" <workers> <max-level>. E.g. `./build/scraper "http://wiprodigital.com/" 3 4`

## Solution

The solution is pretty straightforward: start with the home url, get the content, extract the static links (images,
CSS files, JS scripts), internal and external links; change the urls to absolute, identify the urls that haven't been
scraped before / the urls not in the process of being scraped and add them to the list as long as doing that will take
us to 'deep' in the website (in case we don't want to scrape all pages).

Internally the process starts several GO routines that scrape in parallel and the access to the generated data is
syncronized with a mutex.

The output, besides a few general information, is in JSON format, you can find an example in test_data/result.json.

## What can be done with more time

- Unit testing for everything
- Remove duplicates from results for the same page
- Extract more static data: static CSS files are not scraped, but they can contain links to images, fonts
- Retry failed pages, generate report with the failed requests for scraping at a later time; pause / play functionality
- A better multiple producers / multiple consumers solution using channels, to signal when urls are available
- Increase / Reduce the scraping speed depending on the number of errors
- Add option to scrape through a proxy
- Read last-updated headers and save the data to avoid rescraping stale content
- Do not scrape urls not wanted to be tracked
