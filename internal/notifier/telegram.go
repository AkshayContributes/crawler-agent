package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TelegramClient struct {
	Token  string
	ChatID string
}

func NewTelegramClient(token, chatId string) *TelegramClient {
	return &TelegramClient{
		Token:  token,
		ChatID: chatId,
	}
}

func (tc *TelegramClient) SendMessage(message string) error {
	// Implementation to send message via Telegram Bot API
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tc.Token)

	const maxLen = 4000

	if len(message) > maxLen {
		// Option A: Truncate (Simpler)
		message = message[:maxLen] + "\n\n... [Truncated] ..."

		// Option B (Advanced): You could split into multiple messages here if you wanted
	}

	payload := map[string]string{
		"chat_id": tc.ChatID,
		"text":    message,
	}

	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Printf("Error sending Telegram message: %v\n", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Non-OK HTTP status: %s\n", resp.Status)
		return fmt.Errorf("failed to send message, status: %s", resp.Status)
	}

	return nil

}
