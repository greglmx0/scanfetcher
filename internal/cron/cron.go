package cron

import (
	"fmt"
	"log"
	"time"

	"scanfetcher/internal/domain"
	"scanfetcher/internal/telegram"

	"github.com/robfig/cron"
)

type APIResponse struct {
	Result interface{} `json:"result"`
	Count  int         `json:"count"`
}

// InitCron initialise et démarre le planificateur de tâches cron
func InitCron(telegramBot *telegram.TelegramBot, telegramChatID int64) {
	c := cron.New()
	c.AddFunc("@every 7s", func() {
		log.Println("Exécution de la tâche cron ", time.Now())
		// websites, err := webSiteRepo.GetAllWebsites()
		// if err != nil {
		// 	log.Printf("Erreur lors de la récupération des sites web: %v", err)
		// 	return
		// }

		// newMissions := fetchAndProcessMissions(websites, missionRepo)

		// // Envoi des nouvelles missions par Telegram
		// sendTelegramMessages(telegramBot, telegramChatID, newMissions)
	})
	c.Start()
}

// func fetchAndProcessMissions(websites []domain.Website, missionRepo *repository.MissionRepository) []domain.Mission {
// 	var newMissions []domain.Mission

// 	for _, website := range websites {
// 		log.Printf("Traitement du site web: %v", website.Name)

// 		switch website.Source {
// 		case "VIE":
// 			rawMissions, err := PostRequest(website.URL, website.Body)
// 			if err != nil {
// 				log.Printf("Erreur lors de la requête POST pour %s: %v", website.Name, err)
// 				continue
// 			}

// 			missions, err := decodeVIEResponse(rawMissions)
// 			if err != nil {
// 				log.Printf("Erreur lors du décodage de la réponse pour %s: %v", website.Name, err)
// 				continue
// 			}

// 			newMissions = append(newMissions, processVIEMissions(missions, website.Name, missionRepo)...)
// 		}
// 	}
// 	return newMissions
// }

func sendTelegramMessages(telegramBot *telegram.TelegramBot, telegramChatID int64, missions []domain.Scan) {
	for _, mission := range missions {
		message := fmt.Sprintf("Nouveau scan trouvé: %s\nURL: %s\nDernière lecture: %s\nID du site web: %d",
			mission.Name, mission.Url, mission.LastScanRead, mission.WebsiteID)
		log.Printf("Envoi du message Telegram: %s", message)

		telegramBot.SendMessage(telegramChatID, message)
	}
}

// func PostRequest(url string, body string) (interface{}, error) {
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
// 	if err != nil {
// 		return nil, fmt.Errorf("erreur lors de la création de la requête POST: %v", err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{Timeout: 10 * time.Second}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("erreur lors de l'envoi de la requête POST: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("statut HTTP inattendu: %s", resp.Status)
// 	}

// 	var apiResponse APIResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
// 		return nil, fmt.Errorf("erreur lors du décodage JSON: %v", err)
// 	}

// 	return apiResponse.Result, nil
// }
