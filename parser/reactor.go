package parser

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"main/db"
	"net/http"
	"strings"
)

const postUrl = "http://joyreactor.cc/post/"

var artTags = []string{"#3Dэротика", "#artбарышня", "#ero_art", "#арт_барышня"}

func Request() {
	res, err := http.Get("http://joyreactor.cc/tag/erotic/new")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var posts []db.Post

	// Find the review items
	doc.Find("#post_list .postContainer").Each(func(i int, postEl *goquery.Selection) {
		var post = db.Post{}

		postId, _ := postEl.Attr("id")
		post.Id = strings.Replace(postId, "postContainer", "", 1)

		postEl.Find(".taglist b a").Each(func(i int, s *goquery.Selection) {
			tag := s.Text()
			if !strings.Contains(strings.ToLower(tag), "эрот") {
				tag = strings.ReplaceAll(tag, " ", "_")
				tag = strings.ReplaceAll(tag, "-", "_")
				tag = "#" + tag
				post.Tags = append(post.Tags, tag)
			}
		})

		if artExists(post.Tags) {
			return
		}

		postEl.Find(".post_content .image").Each(func(i int, imageEl *goquery.Selection) {

			videoEl := imageEl.Find("video source[type='video/mp4']")
			if videoEl.Length() > 0 {
				videoEl.Each(func(i int, srcEl *goquery.Selection) {
					src, _ := srcEl.Attr("src")
					post.Gifs = append(post.Gifs, src)
				})
			} else {
				imageEl.Find("img").Each(func(i int, srcEl *goquery.Selection) {
					src, _ := srcEl.Attr("src")
					post.Images = append(post.Images, src)
				})
			}

		})
		posts = append(posts, post)
	})

	db.InsertPosts(posts)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func artExists(tags []string) bool {
	for _, tag := range tags {
		if contains(artTags, tag) {
			return true
		}
	}
	return false
}
