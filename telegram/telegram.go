package telegram

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	PathSetWebhook     = "/setWebhook"
	PathSendMessage    = "/sendMessage"
	PathSendPhoto      = "/sendPhoto"
	PathSendMediaGroup = "/sendMediaGroup"
)

const API_URL = "https://api.telegram.org/bot"

func BuildUrl(param string) string {
	return API_URL + os.Getenv("BOT_TOKEN") + param
}

func HandleWebhook(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	body, _ := ioutil.ReadAll(req.Body)
	log.Println(string(body))
	var update Update
	err := json.Unmarshal(body, &update)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(update)
	err = handleMessage(update.Message)

	if err != nil {
		log.Fatal(err)
	}
}

func handleMessage(msg Message) error {
	return sendMessage(msg)
}

func SendPhoto(msg Message, srcUrl string, caption string) error {
	// outData := OutPhoto{ChatId: msg.Chat.Id, Photo: srcUrl}
	outPhoto := OutPhoto{
		ChatId:  msg.Chat.Id,
		Photo:   srcUrl,
		Caption: caption,
	}
	outPhoto.ReplayToMsgId = msg.Id
	return sendJson(PathSendPhoto, outPhoto)
}

func sendMessage(msg Message) error {
	outData := OutMessage{
		ChatId:        msg.Chat.Id,
		Text:          ">>" + msg.Text,
		ReplayToMsgId: msg.Id,
	}
	return sendJson(PathSendMessage, outData)
}

func sendJson(urlPath string, outData interface{}) error {
	body, err := json.Marshal(outData)
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = http.Post(BuildUrl(urlPath), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func DeleteWebhook() {
	http.Get(BuildUrl(PathSetWebhook + "?url="))
}

func SetWebhook() {
	http.Get(BuildUrl(PathSetWebhook + "?url=" + os.Getenv("HOST") + os.Getenv("BOT_TOKEN")))
}