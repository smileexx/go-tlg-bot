package main

import (
	"io"
	"log"
	"main/parser"
	"main/telegram"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Println("============== Run ==============")

	go updateInterval()
	// db.SelectPostByTag("#photo")

	if os.Getenv("MODE") == "local" {
		log.Println("Run UpdateLoop")
		// remove webhook for local server
		telegram.DeleteWebhook()
		updateLoop()
	} else {
		log.Println("Run ListenAndServe")
		telegram.SetWebhook()
		listenServer()
	}
}

func listenServer() {
	// Wake Up on cron
	http.HandleFunc("/wakeup", func(w http.ResponseWriter, req *http.Request) {
		telegram.SetWebhook()
		io.WriteString(w, "Awake!")
	})

	// Handle Bot WebHook updates
	http.HandleFunc("/"+os.Getenv("BOT_TOKEN"), handleTelegramWebhook)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func updateInterval() {
	for {
		time.Sleep(30 * time.Minute)
		parser.Request("")
	}
}
