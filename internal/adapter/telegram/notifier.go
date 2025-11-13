package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"strings"
	"time"

	"github.com/gobtronic/steam-purchase-notifier/internal/domain"
)

type TelegramNotifier struct {
	botToken string
	chatID   int64
}

func NewTelegramNotifier(botToken string, chatID int64) *TelegramNotifier {
	return &TelegramNotifier{
		botToken: botToken,
		chatID:   chatID,
	}
}

func (n *TelegramNotifier) Notify(game domain.Game) error {
	text := fmt.Sprintf(`*%s*

*[%s](%s)*`, sanitizeString("💸 A new game is available in your Steam Family library!"), sanitizeString(game.Name), sanitizeString(game.StoreURL))
	return n.notify(text, dictionary{})
}

type dictionary = map[string]any

func (n *TelegramNotifier) notify(text string, additionalInfo dictionary) error {
	payload, err := n.buildPayload(text, additionalInfo)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(payload)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.botToken)
	req, err := http.NewRequest(http.MethodGet, url, bodyReader)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (n *TelegramNotifier) buildPayload(text string, additionalInfo dictionary) ([]byte, error) {
	data := dictionary{
		"chat_id":    n.chatID,
		"text":       text,
		"parse_mode": "MarkdownV2",
	}
	maps.Copy(data, additionalInfo)
	dataJson, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}

	return dataJson, nil
}

func sanitizeString(str string) string {
	chars := []rune{
		'_',
		'*',
		'[',
		']',
		',',
		'(',
		')',
		'~',
		'`',
		'>',
		'#',
		'+',
		'-',
		'=',
		'|',
		'{',
		'}',
		'.',
		'!',
	}
	for _, c := range chars {
		str = strings.ReplaceAll(str, string(c), fmt.Sprintf(`\%s`, string(c)))
	}
	return str
}
