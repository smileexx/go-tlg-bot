package reactor

import (
	"log"
	"main/db"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const PostUrl = "http://joyreactor.cc/post/"

var artTags = []string{"#3Dэротика", "#artбарышня", "#ero_art", "#арт_барышня"}

const (
	MediaTypeImg = "img"
	MediaTypeMp4 = "mp4"
	MediaTypeGif = "gif"
)

const (
	videoSelector   = "video source[type='video/mp4']"
	fullImgSelector = "a.prettyPhotoLink"
	imgSelector     = "img"
)

func Parse(page string) {
	if page != "" {
		page = "/" + page
	}
	url := "http://joyreactor.cc/tag/erotic/new" + page
	log.Println(url)
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
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
			media := db.Media{Shown: false}
			if srcEl := imageEl.Find(videoSelector).First(); srcEl.Nodes != nil {
				media.Src, _ = srcEl.Attr("src")
				media.Type = MediaTypeMp4
			} else if srcEl := imageEl.Find(fullImgSelector).First(); srcEl.Nodes != nil {
				media.Src, _ = srcEl.Attr("href")
				media.Type = MediaTypeImg
			} else if srcEl := imageEl.Find(imgSelector).First(); srcEl.Nodes != nil {
				media.Src, _ = srcEl.Attr("src")
				media.Type = MediaTypeImg
			}
			post.Media = append(post.Media, media)
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
