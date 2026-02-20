package scraper

import (
	"strconv"
	"strings"
	"time"
	"net/http"
	"net/url"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	// "go.uber.org/zap/zapcore"
)

type Scraper interface {
	FetchCurrentPrice(itemURL string) (float64, error)
}

type goQueryScraper struct {
	client *http.Client
	selectors map[string]string
	logger *zap.Logger
}

func NewGoQueryScrapper(logger *zap.Logger) * goQueryScraper{ 
	return &goQueryScraper {
		client: &http.Client{Timeout: 10 * time.Second},
		logger: logger,
		selectors: map[string]string{
			"future-phone.ru": ".sp_price span",
		},
	}
}

func(s *goQueryScraper) FetchCurrentPrice(itemURL string) (float64, error) {

	s.logger.Debug("Attempting to fetch price", zap.String("url", itemURL))

	// parse url
	parsedURL, err := url.Parse(itemURL)
	if err != nil {
		return 0, fmt.Errorf("invalid url: %w", err)
	}

	// extract clean url without www.
	hostname := parsedURL.Hostname()
	hostname = strings.TrimPrefix(hostname, "www.")

	// find appropriate sellector in map
	siteSelector, ok := s.selectors[hostname]
	if !ok {
		return 0, fmt.Errorf("no selectors found for domain: %s", hostname)
	}

	// prepare request for sending
	req, err := http.NewRequest("GET", itemURL, nil)
	if err != nil {
		return 0, err
	}

	// change headers to escape ban from server
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")


	// do sends http request and return http response
	// open network connection
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	// close network connection 
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("site returned error status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to parse html: %w", err)
	}
	

	priceStr := doc.Find(siteSelector).First().Text() 
	if priceStr == "" {
		return 0, fmt.Errorf("price element not found with selector: %s", siteSelector)
	}
	
	// remove all rubbish from price (like $, _, space and etc)
	price := s.parsePrice(priceStr)
	if price == 0 {
		return 0, fmt.Errorf("failed to parse price from string: %s", priceStr)
	}

	return price, nil
}

func (s *goQueryScraper) parsePrice(priceStr string) float64 {

	cleanStr := strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' || r == ',' {
			return r
		}
		return -1 // means delete this character
	}, priceStr)

	cleanStr = strings.ReplaceAll(cleanStr, "," , ".")

	price, err := strconv.ParseFloat(cleanStr, 64)
	if err != nil {
		return 0
	}

	return price
}