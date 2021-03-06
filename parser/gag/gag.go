package gag

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"main/parser"
	"math/rand"
	"net/http"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
)

const url = "https://9gag.com/v1/group-posts/tag/funny/fresh"

func Parse() ([]Post, error) {
	var posts []Post
	update := httpGet(url)
	for _, p := range update.Data.Children {
		var post = Post{
			Id:        p.Id,
			Title:     p.Title,
			Permalink: p.Permalink,
		}
		if p.Type == "Animation" {
			post.Type = parser.MediaTypeMp4
			post.Src = p.Images.Animated.Src
		} else {
			post.Type = parser.MediaTypeImg
			post.Src = p.Images.Image.Src
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func GetRandomPost() (Post, error) {
	var post Post
	posts, err := Parse()
	if err != nil {
		return post, err
	}
	i := rand.Intn(len(posts))
	return posts[i], nil
}

func httpGet(url string) Update {
	var update Update
	client := &http.Client{}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

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
