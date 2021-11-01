package reddit

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const url = "https://www.reddit.com/r/ProgrammerHumor.json"

var sources = []string{
	"https://www.reddit.com/r/ProgrammerHumor.json",
}

func Parse() ([]Post, error) {
	// var updates []Update
	// for _, url := range sources {
	// update := httpGet(url)
	// 	updates = append(updates, update)
	// }
	var posts []Post
	update := httpGet(url)
	for _, p := range update.Data.Children {
		posts = append(posts, Post(p.Data))
	}
	return posts, nil
}

func httpGet(url string) Update {
	var update Update
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", " Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &update)

	return update
}
