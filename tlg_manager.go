package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/db"
	"main/parser"
	"main/telegram"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const postBatchCount = 30

var postsBuffer []db.Post

func init() {
	db.CreateConnection()
}

func handleTelegramWebhook(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	body, _ := ioutil.ReadAll(req.Body)
	log.Println(string(body))
	var update telegram.Update
	err := json.Unmarshal(body, &update)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(update)
	err = reactOnMessage(update.Message)

	if err != nil {
		log.Fatal(err)
	}
}

func reactOnMessage(msg telegram.Message) error {
	//TODO: do something if needs

	if strings.HasPrefix(msg.Text, "/boobs") {
		return commandBoobs(msg)
	}

	if strings.HasPrefix(msg.Text, "/tag") {
		return commandTag(msg)
	}

	return nil
}

func commandBoobs(msg telegram.Message) error {
	log.Println("buffer len", len(postsBuffer))
	if len(postsBuffer) < 1 {
		postsBuffer = db.GetRandomPosts(postBatchCount)
	}

	// exit if no posts
	if len(postsBuffer) < 1 {
		return nil
	}

	i := rand.Intn(len(postsBuffer))
	post := postsBuffer[i]
	media := post.Media
	// remove used item
	postsBuffer = append(postsBuffer[:i], postsBuffer[i+1:]...)

	j := rand.Intn(len(media))
	item := media[j]
	caption := strings.Join(post.Tags, " ") + "\n" + parser.PostUrl + post.Id
	if item.Type == parser.MediaTypeImg {
		return telegram.SendPhoto(msg, item.Src, caption)
	} else {
		return telegram.SendVideo(msg, item.Src, caption)
	}
}

func commandTag(msg telegram.Message) error {
	regex := *regexp.MustCompile(`#[_\wА-Яа-я]+`)
	tag := regex.FindStringSubmatch(msg.Text)
	if len(tag) < 1 {
		return nil
	}
	posts := db.SelectPostsByTag(tag[0])
	// exit if no posts
	if len(posts) < 1 {
		return nil
	}
	i := rand.Intn(len(posts))
	post := posts[i]
	media := post.Media
	j := rand.Intn(len(media))
	item := media[j]
	caption := strings.Join(post.Tags, " ") + "\n" + parser.PostUrl + post.Id
	if item.Type == parser.MediaTypeImg {
		return telegram.SendPhoto(msg, item.Src, caption)
	} else {
		return telegram.SendVideo(msg, item.Src, caption)
	}
}

/**
| ==================== For local =====================
*/
func getUpdates(offset int) ([]telegram.Update, error) {
	// http.NewRequest("GET", API_URL+"/getMe")
	resp, err := http.Get(telegram.BuildUrl("/getUpdates?offset=" + fmt.Sprint(offset)))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))

	var restResponse telegram.RestResponse
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
		updates, err := getUpdates(offset)
		if err != nil {
			log.Fatal(err)
		}
		for _, update := range updates {
			offset = update.UpdateId + 1
			log.Println(update.Message)
			err = reactOnMessage(update.Message)
			if err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(time.Second)
	}
}
