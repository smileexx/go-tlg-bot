package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const API_URL = "https://api.telegram.org/bot"

func main() {
	log.Println("============== Run ==============")

	// Wake Up on cron
	http.HandleFunc("/wakeup", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Awake!")
	})

	// Hnadle Bot Webhook updates
	http.HandleFunc("/"+os.Getenv("BOT_TOKEN"), func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		body, _ := ioutil.ReadAll(req.Body)
		fmt.Printf("%s", body)
	})
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	// offset := 0

	// for {
	// 	currentTime := time.Now()
	// 	log.Println(currentTime.Format("2006-01-02 15:04:05.000000"))
	// 	updates, err := getUpdates(offset)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	for _, update := range updates {
	// 		offset = update.UpdateId + 1
	// 		log.Println(update.Message)
	// 		err = respond(update.Message)
	// 	}
	// 	time.Sleep(time.Second)
	// }

}

func buildUrl(param string) string {
	return API_URL + os.Getenv("BOT_TOKEN") + param
}

func getUpdates(offset int) ([]Update, error) {
	// http.NewRequest("GET", API_URL+"/getMe")
	resp, err := http.Get(buildUrl("/getUpdates?offset=" + fmt.Sprint(offset)))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s", body)

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return restResponse.Result, nil
}

func respond(msg Message) error {
	botMsg := BotMessage{ChatId: msg.Chat.ChatId, Text: ">>" + msg.Text}
	body, err := json.Marshal(botMsg)
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = http.Post(buildUrl("/sendMessage"), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
