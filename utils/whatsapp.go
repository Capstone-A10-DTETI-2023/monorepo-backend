package utils

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var tokenWA = os.Getenv("TOKEN_WA")

var (
	notifCooldownEnd = make(map[string]int64)
)

func SendWAMessage(phoneNum, message, schedule string) error {

	if time.Now().Unix() < notifCooldownEnd[phoneNum] {
		return fiber.ErrTooManyRequests
	}

	if tokenWA == "" {
		return fiber.ErrInternalServerError
	}
	
	if schedule == "" {
		schedule = "0"
	}

	data := url.Values{}
	data.Set("target", phoneNum)
	data.Set("message", message)
	data.Set("schedule", schedule)

	notifCooldownEnd[phoneNum] = time.Now().Add(5 * time.Minute).Unix()

	client := &http.Client{}
	r, err := http.NewRequest(http.MethodPost, "https://api.fonnte.com/send", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Authorization", tokenWA)

	result, err := client.Do(r)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	return nil
}