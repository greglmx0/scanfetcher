package notifier

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot *tgbotapi.BotAPI
}

type TelegramNotifier struct {
	bot    *TelegramBot
	chatID int64
}

func NewTelegramNotifier(bot *TelegramBot, chatID int64) *TelegramNotifier {
	return &TelegramNotifier{bot: bot, chatID: chatID}
}

func NewTelegramBot(token string) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	return &TelegramBot{bot: bot}, nil
}

func (tb *TelegramBot) SendMessage(chatID int64, message string) error {

	sendMessageTelegram := os.Getenv("SEND_MESSAGE_TELEGRAM")
	if sendMessageTelegram == "" {
		sendMessageTelegram = "true"
	}

	if sendMessageTelegram == "false" {
		log.Println(`La variable d'environnement "SEND_MESSAGE_TELEGRAM" est définie à "false", le message [` + message + `] ne sera pas envoyé`)
		return nil
	} else {

		msg := tgbotapi.NewMessage(chatID, message)
		_, err := tb.bot.Send(msg)

		if err != nil {
			log.Fatalf("Erreur lors de l'envoi du message: %v", err)
			return err
		} else {
			log.Printf("Message envoyé: %s", message)
		}
	}
	return nil
}
