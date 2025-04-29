package scraper

import "fmt"

func GetScraper(siteName string) (Scraper, error) {
	switch siteName {
	case "mangas-origines.fr":
		return MangasOriginesScraper{}, nil
	default:
		return nil, fmt.Errorf("aucun scraper défini pour %s", siteName)
	}
}
