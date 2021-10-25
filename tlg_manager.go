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
	if strings.HasPrefix(msg.Text, "/post") {
		return commandPost(msg)
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
	// remove used item
	postsBuffer = append(postsBuffer[:i], postsBuffer[i+1:]...)

	return sendSingleMedia(msg, post)
}

func commandTag(msg telegram.Message) error {
	regex := *regexp.MustCompile(`#[_\wА-Яа-я]+`)
	tags := regex.FindStringSubmatch(msg.Text)
	// exit if no tag
	if len(tags) < 1 {
		telegram.SendMessage(msg, "Command should contains > #tag")
		return nil
	}
	tag := tags[0]
	posts, err := db.SelectPostsByTag(tag)
	// exit if no posts
	if err != nil {
		telegram.SendMessage(msg, "Nothing found by tag "+tag)
		return err
	}
	i := rand.Intn(len(posts))
	post := posts[i]
	return sendSingleMedia(msg, post)
}

func commandPost(msg telegram.Message) error {
	regex := *regexp.MustCompile(`/post\s+(\d+)`)
	submatch := regex.FindStringSubmatch(msg.Text)
	// exit if no tag
	if len(submatch) < 2 {
		telegram.SendMessage(msg, "Command should be like > `/post 4958161`")
		return nil
	}
	postId := submatch[1]
	post, err := db.SelectPostsById(postId)
	// exit if no posts
	if err != nil {
		telegram.SendMessage(msg, "Nothing found by post ID "+postId)
		return nil
	}
	if len(post.Media) > 1 {
		return sendGroupedMedia(msg, *post)
	} else if len(post.Media) == 1 {
		return sendSingleMedia(msg, *post)
	}

	return nil
}

func sendSingleMedia(msg telegram.Message, post db.Post) error {
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

func sendGroupedMedia(msg telegram.Message, post db.Post) error {
	var mediaItems []telegram.InputMediaItem
	var err error
	for _, media := range post.Media {
		input := telegram.InputMediaItem{Media: media.Src}
		if media.Type == parser.MediaTypeImg {
			input.Type = "photo"
		} else if media.Type == parser.MediaTypeMp4 {
			err = telegram.SendVideo(msg, media.Src, "")
			continue
		}
		if len(mediaItems) == 0 {
			input.Caption = strings.Join(post.Tags, " ") + "\n" + parser.PostUrl + post.Id
		}
		mediaItems = append(mediaItems, input)
	}
	if len(mediaItems) > 0 {
		err = telegram.SendMediaGroup(msg, mediaItems)
	}
	return err
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
