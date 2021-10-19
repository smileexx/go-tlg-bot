package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const BOT_TOKEN = "2036907836:AAGIRjFvj4HGFRKZ9kBysjrMPXJ6oOpZ278"
const API_URL = "https://api.telegram.org/bot" + BOT_TOKEN

func main() {
	getUpdates()
}

func getUpdates() ([]Update, error) {
	// http.NewRequest("GET", API_URL+"/getMe")
	resp, err := http.Get(API_URL + "/getUpdates")
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
