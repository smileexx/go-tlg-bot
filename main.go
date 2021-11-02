package main

import (
	"io"
	"log"
	"main/parser/reactor"
	"main/parser/reddit"
	"main/telegram"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Println("============== Run ==============")

	go updateInterval()
	go onTimerEvents()
	// db.SelectPostByTag("#photo")

	telegram.SetCommands()

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

/**
 *	Parse content updates from
 */
func updateInterval() {
	for {
		reactor.Parse("")
		reddit.Parse()
		time.Sleep(30 * time.Minute)
	}
}

/**
 *	Send news for subscribers
 */
func onTimerEvents() {
	for {
		now := time.Now()
		nowTs := now.Unix()
		schedule(nowTs)

		time.Sleep(10 * time.Minute)
	}
}
