package service

import (
	"fmt"
	"net/http"
	"price_checker/internal/core/domains"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Repository interface {
	Add(domains.Item) (domains.Item, error)
	Delete(int64) error
	GetAll() ([]domains.Item, error)
	UpdatePrice(int64, float64) error
}

type PriceService struct {
	repo Repository
}

func NewPriceService(repo Repository) *PriceService {
	return &PriceService{repo: repo}
}

func (s *PriceService) CreateItem(item domains.Item) (domains.Item, error) {

	currentPrice, err := s.fetchCurrentPrice(item.URL)
	if err != nil {
		return domains.Item{}, fmt.Errorf("couldn't fetch price: %w", err)
	}

	if currentPrice == 0 {
		return domains.Item{}, fmt.Errorf("failed to parse price from the page")
	}
	
	item.CurrentPrice = currentPrice	
	
	return s.repo.Add(item)
}

func (s *PriceService) ListItems() ([]domains.Item, error) {
	return s.repo.GetAll()
}

func(s *PriceService) fetchCurrentPrice(itemURL string) (float64, error) {

	// skip after 10 sec if site is'n working
	client := &http.Client{Timeout: 10 * time.Second}

	// sent request on particular URL
	req, err := http.NewRequest("GET", itemURL, nil)
	if err != nil {
		return 0, err
	}

	// set header User-Agent (показывает откуда пришёл запрос - если бы этой строки небыло, то postman)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")	

	// do sends http request and return http response
	// open network connection
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	// close network connection 
	defer resp.Body.Close()

	fmt.Println("Status code", resp.StatusCode)
	// get http response from the other hand - if not 200 --> error
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("failed to fetch page: status code %d", resp.StatusCode)
	}

	// goquery library convert html text to DOM tree
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return 0, err 
	}

	// first - return first tag from all found tags. Text - extract text from this tag
	priceStr := doc.Find(".sp_price span").First().Text()
	
	// remove all rubbish from price (like $, _, space and etc)
	price := s.parsePrice(priceStr)

	return price, nil
}

func (s *PriceService) parsePrice(priceStr string) float64 {

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