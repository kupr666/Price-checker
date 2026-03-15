package scraper

import (
	"context"

	"price_checker/internal/features/price_tracker/cache"

	"go.uber.org/zap"
)

type CachedScraper struct {
	scraper Scraper
	cache cache.PriceCache
	logger *zap.Logger
}

func NewCacheScraper(scraper Scraper, cache cache.PriceCache, logger *zap.Logger) *CachedScraper{
	return&CachedScraper{
		scraper: scraper,
		cache: cache,
		logger: logger,
	}
}

func (c *CachedScraper) FetchCurrentPrice(ctx context.Context, itemURL string) (float64, error) {

	// get error
	price, found, err := c.cache.Get(ctx, itemURL)
	if err != nil {
		c.logger.Warn("redis cache get failed, falling back to scraper", zap.Error(err))
	}

	// found cache
	if found {
		c.logger.Debug("cache found", zap.String("url:", itemURL), zap.Float64("price", price))
		return price, nil
	}

	c.logger.Debug("cache miss, scraping", zap.String("url", itemURL))

	// if the cache is empty, we fetch the price from the website
	price, err = c.scraper.FetchCurrentPrice(ctx, itemURL)
	if err != nil {
		return 0, err
	}

	//
	if err := c.cache.Set(ctx, itemURL, price); err != nil {
		c.logger.Warn("redis cache set failed", zap.Error(err), zap.String("url:", itemURL))
	}

	return price, nil

}