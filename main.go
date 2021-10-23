package main

import (
	"io"
	"log"
	"main/telegram"
	"net/http"
	"os"
)

func main() {
	log.Println("============== Run ==============")

	if os.Getenv("MODE") == "local" {
		log.Println("Run UpdateLoop")
		// remove webhook for local server
		telegram.DeleteWebhook()
		telegram.UpdateLoop()
	} else {
		log.Println("Run ListenAndServe")
		telegram.SetWebhook()
		listenServer()
	}
}

func listenServer() {
	// Wake Up on cron
	http.HandleFunc("/wakeup", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Awake!")
	})

	// Handle Bot WebHook updates
	http.HandleFunc("/"+os.Getenv("BOT_TOKEN"), telegram.HandleWebhook)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
