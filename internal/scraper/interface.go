package scraper

import (
	"scanfetcher/internal/domain"

	"github.com/go-rod/rod"
)

type Scraper interface {
	Scrape(browser *rod.Browser, url string) ([]domain.Scan, error)
}

type Builder func() Scraper

type Factory interface {
	Register(name string, b Builder)
	Get(name string) (Scraper, error)
}
