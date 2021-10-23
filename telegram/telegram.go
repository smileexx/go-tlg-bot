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
	PathSendMessage    = "/sendMessage"
	PathSendPhoto      = "/sendPhoto"
	PathSendMediaGroup = "/sendMediaGroup"
)

const API_URL = "https://api.telegram.org/bot"

func buildUrl(param string) string {
	return API_URL + os.Getenv("BOT_TOKEN") + param
}

func HandleUpdate(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	body, _ := ioutil.ReadAll(req.Body)
	log.Println(string(body))
	var update Update
	err := json.Unmarshal(body, &update)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(update)
	err = sendMessage(update.Message)

	if err != nil {
		log.Fatal(err)
	}
}

func sendPhoto(msg Message) error {
	outData := OutPhoto{ChatId: msg.Chat.Id, Photo: "http://"}
	return sendJson(PathSendMessage, outData)
}

func sendMessage(msg Message) error {
	outData := OutMessage{ChatId: msg.Chat.Id, Text: ">>" + msg.Text}
	return sendJson(PathSendMessage, outData)
}

func sendJson(urlPath string, outData interface{}) error {
	body, err := json.Marshal(outData)
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = http.Post(buildUrl(urlPath), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

/*
func getUpdates(offset int) ([]Update, error) {
	// http.NewRequest("GET", API_URL+"/getMe")
	resp, err := http.Get(buildUrl("/getUpdates?offset=" + fmt.Sprint(offset)))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return restResponse.Result, nil
}

func updateLoop() {

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
			if err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(time.Second)
	}
}
*/
