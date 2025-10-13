package notifier

type Notifier interface {
	SendMessage(chatID int64, text string) error
}
