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

	//db.Connect()
	//os.Exit(0)

	// Wake Up on cron
	http.HandleFunc("/wakeup", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Awake!")
	})

	// Handle Bot WebHook updates
	http.HandleFunc("/"+os.Getenv("BOT_TOKEN"), telegram.HandleUpdate)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	// updateLoop()
}
