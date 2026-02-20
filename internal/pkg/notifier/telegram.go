package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"
)

type TelegramNotifier struct {
	token string
	chatID string
}

func NewTelegramNotifier(token, chatID string) *TelegramNotifier {
	return &TelegramNotifier{token: token, chatID: chatID}
}

func (t *TelegramNotifier) Notify(message string) error {

	// get the url with particular token
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)

	
	payload := map[string]string{
		"chat_id": t.chatID,
		"text":    message,
	}

	// json marshal returns a slice of bytes in body
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url,"application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api error: status %d", resp.StatusCode)
	}

	return nil
}