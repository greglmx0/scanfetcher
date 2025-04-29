package cron

import (
	"fmt"
	"log"
	"os"
	"time"

	"scanfetcher/internal/domain"
	"scanfetcher/internal/repository"
	"scanfetcher/internal/scraper"
	"scanfetcher/internal/telegram"

	"github.com/go-rod/rod"
	"github.com/robfig/cron"
)

type APIResponse struct {
	Result interface{} `json:"result"`
	Count  int         `json:"count"`
}

// InitCron initialise et démarre le planificateur de tâches cron
func InitCron(telegramBot *telegram.TelegramBot, telegramChatID int64, webSiteRepo *repository.WebsiteRepository, scanRepo *repository.ScanRepository) {

	cr := cron.New()
	cron_job_schedule := os.Getenv("CRON_JOB_SCHEDULE")
	log.Printf("Planificateur de tâches cron démarré avec la planification: %s", cron_job_schedule)
	cr.AddFunc(cron_job_schedule, func() {
		log.Println("Exécution de la tâche cron ", time.Now())
		websites, err := webSiteRepo.GetAll()
		if err != nil {
			log.Printf("Erreur lors de la récupération des sites web: %v", err)
			return
		}

		log.Printf("Sites web récupérés: %v", websites)
		newScans := fetchAndProcessMissions(websites, scanRepo)

		sendTelegramMessages(telegramBot, telegramChatID, newScans)
	})
	cr.Start()
}

func fetchAndProcessMissions(websites []domain.Website, scanRepo *repository.ScanRepository) []domain.Scan {
	var allScans []domain.Scan

	for _, site := range websites {

		scan, _ := scanRepo.GetByWebsiteID(site.ID)

		fmt.Printf("Site: %s, Scans: %v\n", site.Name, scan)

		for _, s := range scan {

			browser := rod.New().MustConnect()
			scraper, err := scraper.GetScraper(site.Name)
			if err != nil {
				log.Printf("Scraper non disponible pour %s: %v", site.Name, err)
				continue
			}

			log.Printf("Lancement du scraper pour %s", site.Name)
			scans, err := scraper.Scrape(browser, s.Url)
			browser.MustClose()
			if err != nil {
				log.Printf("Erreur lors du scraping de %s: %v", site.Name, err)
				continue
			}

			// check if the last scan read is the same as the last scan
			if len(scans) == 0 || scans[0].LastScanRead == s.LastScanRead {
				log.Printf("Aucun nouveau scan trouvé pour %s", site.Name)
				continue
			}

			// On met à jour le dernier scan lu
			err = scanRepo.UpdateLastScanRead(int(s.ID), scans[0].LastScanRead)
			if err != nil {
				log.Printf("Erreur lors de la mise à jour du dernier scan lu pour %s: %v", site.Name, err)
				continue
			}
			log.Printf("Dernier scan lu mis à jour pour %s", site.Name)

			// On envoie le message Telegram
			message := fmt.Sprintf("Nouveau scan trouvé: %s\nURL: %s\nDernière lecture: %s\nID du site web: %d",
				scans[0].Name, scans[0].Url, scans[0].LastScanRead, s.WebsiteID)
			log.Printf("Envoi du message Telegram: %s", message)

			allScans = append(allScans, scans...)
		}
	}

	return allScans
}

func sendTelegramMessages(telegramBot *telegram.TelegramBot, telegramChatID int64, scans []domain.Scan) {
	for _, s := range scans {
		message := fmt.Sprintf("Nouveau scan trouvé: %s\nURL: %s\nDernière lecture: %s\nID du site web: %d",
			scans[0].Name, scans[0].Url, scans[0].LastScanRead, s.WebsiteID)
		log.Printf("Envoi du message Telegram: %s", message)

		telegramBot.SendMessage(telegramChatID, message)
	}
}
