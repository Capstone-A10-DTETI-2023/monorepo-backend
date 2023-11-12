package utils

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pusher/pusher-http-go/v5"
)

func ConnectWS() *pusher.Client {
	pusherClient := pusher.Client{
		AppID: os.Getenv("PUSHER_APP_ID"),
		Key: os.Getenv("PUSHER_APP_KEY"),
		Secret: os.Getenv("PUSHER_APP_SECRET"),
		Cluster: os.Getenv("PUSHER_APP_CLUSTER"),
		Secure: true,
	}

	log.Println("Connected to Pusher")

	return &pusherClient
}
