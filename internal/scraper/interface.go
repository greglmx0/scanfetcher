package scraper

import (
	"scanfetcher/internal/domain"

	"github.com/go-rod/rod"
)

type Scraper interface {
	Scrape(browser *rod.Browser, url string) ([]domain.Scan, error)
}
