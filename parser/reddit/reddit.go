package reddit

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"main/db"
	"main/parser"
	"net/http"
	"strings"
)

type Update struct {
	Kind string `json:"kind"`
	Data struct {
		After    string `json:"after"`
		Children []struct {
			Data struct {
				Id        string `json:"id"`
				Title     string `json:"title"`
				Type      string
				Src       string  `json:"url_overridden_by_dest"`
				Permalink string  `json:"permalink"`
				Tag       string  `json:"subreddit"`
				Created   float32 `json:"created_utc"`
				IsVideo   bool    `json:"is_video"`
				PostHint  string  `json:"post_hint"`
				Media     struct {
					Video struct {
						Src string `json:"fallback_url"`
					} `json:"reddit_video"`
				} `json:"media"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

const (
	postHintVideo = "hosted:video"
	postHintImage = "image"
	postHintLink  = "link"
)

const host = "https://www.reddit.com"

// const url = "https://www.reddit.com/r/ProgrammerHumor.json"

var sources = []string{
	"https://www.reddit.com/r/ProgrammerHumor/top.json",
	"https://www.reddit.com/r/WTF/new.json",
	"https://www.reddit.com/r/memes/top.json",
}

func Parse() error {
	var posts []db.MemePost
	for _, url := range sources {
		update := httpGet(url)
		for _, p := range update.Data.Children {
			if p.Data.PostHint == postHintLink {
				// skip post if content Src is remote link
				continue
			}
			var post = db.MemePost{
				Id:        p.Data.Id,
				Title:     p.Data.Title,
				Type:      parser.MediaTypeImg,
				Src:       p.Data.Src,
				Permalink: p.Data.Permalink,
				Shown:     false,
				Created:   int(p.Data.Created),
				Tags:      []string{"#" + p.Data.Tag},
			}
			post.Permalink = host + post.Permalink
			if p.Data.IsVideo {
				post.Type = parser.MediaTypeMp4
				post.Src = p.Data.Media.Video.Src
			}
			if strings.HasSuffix(post.Src, parser.MediaTypeGif) {
				post.Type = parser.MediaTypeGif
			}
			if post.Src == "" {
				continue
			}

			posts = append(posts, post)
		}

	}

	return db.InsertMemesPosts(posts)
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
