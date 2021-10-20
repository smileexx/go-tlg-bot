package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const API_URL = "https://api.telegram.org/bot"

func main() {
	fmt.Println("BOT_TOKEN", os.Getenv("BOT_TOKEN"))
	offset := 0

	for {
		currentTime := time.Now()
		log.Println(currentTime.Format("2006-01-02 15:04:05.000000"))
		updates, err := getUpdates(offset)
		if err != nil {
			log.Fatal(err)
		}
		for _, update := range updates {
			offset = update.UpdateId + 1
			log.Println(update.Message)
			err = respond(update.Message)
		}
		time.Sleep(time.Second)
	}
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
	_, err = http.Post(API_URL+"/sendMessage", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
