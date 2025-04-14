package cron

import (
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"scanfetcher/internal/domain"
	"scanfetcher/internal/repository"
	"scanfetcher/internal/telegram"

	"github.com/go-rod/rod"
	"github.com/robfig/cron"
)

type APIResponse struct {
	Result interface{} `json:"result"`
	Count  int         `json:"count"`
}

// InitCron initialise et démarre le planificateur de tâches cron
func InitCron(telegramBot *telegram.TelegramBot, telegramChatID int64, webSiteRepo *repository.WebsiteRepository) {
	// os.Setenv("ROD_BROWSER_PATH", "/usr/bin/chromium")
	// fmt.Println("ROD_BROWSER_PATH:", os.Getenv("ROD_BROWSER_PATH"))
	rob := rod.New().MustConnect()

	fmt.Println("Initialisation du collecteur rob")
	scrapScans(rob, "https://mangas-origines.fr/oeuvre/the-beginning-after-the-end/")
	fmt.Println("Collecte terminée merci rob")

	cr := cron.New()
	cr.AddFunc("@every 10m", func() {
		log.Println("Exécution de la tâche cron ", time.Now())
		websites, err := webSiteRepo.GetAll()
		if err != nil {
			log.Printf("Erreur lors de la récupération des sites web: %v", err)
			return
		}

		log.Printf("Sites web récupérés: %v", websites)

		// newMissions := fetchAndProcessMissions(websites, missionRepo)

		// // Envoi des nouvelles missions par Telegram
		// sendTelegramMessages(telegramBot, telegramChatID, newMissions)
	})
	cr.Start()
}

func fetchAndProcessMissions(websites []domain.Website) []domain.Scan {
	var rowScan []domain.Scan

	for _, website := range websites {
		log.Printf("Traitement du site web: %v", website.Name)

		switch website.Name {
		case "mangas-origines.fr":

		}
		return rowScan
	}
	return rowScan
}

func sendTelegramMessages(telegramBot *telegram.TelegramBot, telegramChatID int64, missions []domain.Scan) {
	for _, mission := range missions {
		message := fmt.Sprintf("Nouveau scan trouvé: %s\nURL: %s\nDernière lecture: %s\nID du site web: %d",
			mission.Name, mission.Url, mission.LastScanRead, mission.WebsiteID)
		log.Printf("Envoi du message Telegram: %s", message)

		telegramBot.SendMessage(telegramChatID, message)
	}
}

func scrapScans(rod *rod.Browser, url string) string {
	page := rod.MustPage(url)
	page.MustWaitLoad()

	// On chope la div principale
	div := page.MustElement("#manga-chapters-holder")

	// Tous les chapitres
	chapters := div.MustElements("li.wp-manga-chapter")
	if len(chapters) == 0 {
		return "Aucun chapitre trouvé"
	}

	// Ne garder que le premier chapitre si plus d'un chapitre est trouvé
	if len(chapters) > 1 {
		chapters = chapters[:1]
	}

	// Reverse les chapitres
	slices.Reverse(chapters)

	var result strings.Builder
	for _, chapter := range chapters {
		a := chapter.MustElement("a")
		title := strings.TrimSpace(a.MustText())
		if title == "" {
			continue
		}

		link := a.MustProperty("href").String()
		if link == "" {
			continue
		}

		date := getDate(chapter)
		views := getViews(chapter)

		result.WriteString(fmt.Sprintf("\n------------\n📘 %s\n🔗 %s\n📅 %s\n👁️  %s\n------------\n", title, link, date, views))
	}

	log.Printf("Résultat de la collecte: %s", result.String())
	return result.String()
}

func getDate(chapter *rod.Element) string {
	dateSpan := chapter.MustElement("span.chapter-release-date")
	if dateSpan == nil {
		return ""
	}

	// Essaye de trouver l'élément <i> dans le span
	if iElement, err := dateSpan.Element("i"); err == nil && iElement != nil {
		return strings.TrimSpace(iElement.MustText())
	}

	// Sinon cherche un <a> dans le span
	if aTag, err := dateSpan.Element("a"); err == nil && aTag != nil {
		return aTag.MustProperty("title").String()
	}

	return ""
}

func getViews(chapter *rod.Element) string {
	viewSpan := chapter.MustElement("span.view")
	if viewSpan != nil {
		return strings.TrimSpace(viewSpan.MustText())
	}
	return ""
}
