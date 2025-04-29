package main

import (
	"log"
	"net/http"
	"os"
	"scanfetcher/internal/cron"
	"scanfetcher/internal/db"

	"scanfetcher/internal/fixtures"
	"strconv"

	handlers "scanfetcher/internal/handler"
	"scanfetcher/internal/repository"
	"scanfetcher/internal/telegram"
	"scanfetcher/internal/usecase"

	"github.com/gorilla/mux"
)

func main() {
	// init telegram bot
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	telegramChatID, err := strconv.ParseInt(telegramChatIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Erreur lors de la conversion de TELEGRAM_CHAT_ID: %v", err)
	}

	log.Println("Initialisation du bot Telegram")
	telegramBot, err := telegram.NewTelegramBot(telegramToken)
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation du bot Telegram: %v", err)
	}

	// Initialiser la base de données
	dbPath := "/app/data"
	dbName := "scanfetcher.db"
	db, err := db.InitDB(dbPath, dbName)
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de la base de données: %v", err)
	}

	fixtures.SeedData(db)

	// Configurer le routeur HTTP
	r := mux.NewRouter()

	// User
	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handlers.NewUserHandler(userUseCase)
	r.HandleFunc("/user/{id:[0-9]+}", userHandler.GetUserHandler).Methods("GET")
	r.HandleFunc("/users", userHandler.GetAllUsersHandler).Methods("GET")
	r.HandleFunc("/user", userHandler.CreateUserHandler).Methods("POST")

	// Website
	websiteRepo := repository.NewWebsiteRepository(db)
	websiteUseCase := usecase.NewWebsiteUseCase(websiteRepo)
	websiteHandlers := handlers.NewWebsiteHandler(websiteUseCase)
	r.HandleFunc("/websites", websiteHandlers.CreateWebsite).Methods("POST")
	r.HandleFunc("/websites", websiteHandlers.ListWebsites).Methods("GET")
	r.HandleFunc("/websites/{id}", websiteHandlers.GetWebsite).Methods("GET")

	// Scan
	scanRepo := repository.NewScanRepository(db)
	scanUC := usecase.NewScanUseCase(scanRepo)
	scanHandler := handlers.NewScanHandler(scanUC)
	r.HandleFunc("/scans", scanHandler.CreateScan).Methods("POST")
	r.HandleFunc("/websites/{websiteID}/scans", scanHandler.GetScansForWebsite).Methods("GET")

	// Initialiser le planificateur de tâches cron
	cron.InitCron(telegramBot, telegramChatID, websiteRepo, scanRepo)

	log.Println("Serveur en cours d'exécution sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
