package scraper

import (
	"fmt"
	"scanfetcher/internal/domain"
	"slices"
	"strings"

	"github.com/go-rod/rod"
)

type MangasOriginesScraper struct{}

func (s MangasOriginesScraper) Scrape(browser *rod.Browser, url string) ([]domain.Scan, error) {
	page := browser.MustPage(url)
	page.MustWaitLoad()

	var scans []domain.Scan

	div := page.MustElement("#manga-chapters-holder")
	chapters := div.MustElements("li.wp-manga-chapter")

	// print the number of chapters found
	fmt.Printf("Found %d chapters\n", len(chapters))

	if len(chapters) == 0 {
		return nil, nil
	}
	if len(chapters) > 1 {
		chapters = chapters[:1]
	}
	slices.Reverse(chapters)

	for _, chapter := range chapters {
		LastScanRead := strings.TrimSpace(strings.Split(chapter.MustText(), " ")[1])
		a := chapter.MustElement("a")
		title := strings.TrimSpace(a.MustText())
		link := a.MustProperty("href").String()
		if title == "" || link == "" {
			continue
		}

		scans = append(scans, domain.Scan{
			Name:         title,
			Url:          link,
			LastScanRead: LastScanRead,
		})
	}

	return scans, nil
}

// func scrapScans(rod *rod.Browser, url string) string {
// 	page := rod.MustPage(url)
// 	page.MustWaitLoad()

// 	// On chope la div principale
// 	div := page.MustElement("#manga-chapters-holder")

// 	// Tous les chapitres
// 	chapters := div.MustElements("li.wp-manga-chapter")
// 	if len(chapters) == 0 {
// 		return "Aucun chapitre trouvé"
// 	}

// 	// Ne garder que le premier chapitre si plus d'un chapitre est trouvé
// 	if len(chapters) > 1 {
// 		chapters = chapters[:1]
// 	}

// 	// Reverse les chapitres
// 	slices.Reverse(chapters)

// 	var result strings.Builder
// 	for _, chapter := range chapters {
// 		a := chapter.MustElement("a")
// 		title := strings.TrimSpace(a.MustText())
// 		if title == "" {
// 			continue
// 		}

// 		link := a.MustProperty("href").String()
// 		if link == "" {
// 			continue
// 		}

// 		date := getDate(chapter)
// 		views := getViews(chapter)

// 		result.WriteString(fmt.Sprintf("\n------------\n📘 %s\n🔗 %s\n📅 %s\n👁️  %s\n------------\n", title, link, date, views))
// 	}

// 	log.Printf("Résultat de la collecte: %s", result.String())
// 	return result.String()
// }

// func getDate(chapter *rod.Element) string {
// 	dateSpan := chapter.MustElement("span.chapter-release-date")
// 	if dateSpan == nil {
// 		return ""
// 	}

// 	// Essaye de trouver l'élément <i> dans le span
// 	if iElement, err := dateSpan.Element("i"); err == nil && iElement != nil {
// 		return strings.TrimSpace(iElement.MustText())
// 	}

// 	// Sinon cherche un <a> dans le span
// 	if aTag, err := dateSpan.Element("a"); err == nil && aTag != nil {
// 		return aTag.MustProperty("title").String()
// 	}

// 	return ""
// }

// func getViews(chapter *rod.Element) string {
// 	viewSpan := chapter.MustElement("span.view")
// 	if viewSpan != nil {
// 		return strings.TrimSpace(viewSpan.MustText())
// 	}
// 	return ""
// }
