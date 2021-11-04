package main

import (
	"main/db"
	"main/telegram"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const MinPeriod = 30 // minutes

const (
	FeedBoobs = "boobs"
	FeedMemes = "memes"
)

var feedsDescription = map[string]string{
	"boobs": "Subscribe to regular ero images",
	"memes": "Subscribe to random memes",
}

func schedule(nowTs int64) error {
	subscriptions, err := db.SelectSchedule()
	// select same content for all subscribers
	var tmp = make(map[string]interface{})
	for _, sub := range subscriptions {
		if nowTs-sub.LastAction >= int64(sub.Period*60) {
			switch sub.Type {
			case "boobs":
				var post db.Post
				if _, ok := tmp[sub.Type]; ok {
					post = tmp[sub.Type].(db.Post)
				} else {
					rPost, err := getRandomPost()
					if err != nil {
						return err
					}
					tmp[sub.Type] = rPost
					post = rPost
				}
				err = sendSingleBoobs(sub.ChatId, post)
				break
			case "memes":
				var meme db.MemePost
				if _, ok := tmp[sub.Type]; ok {
					meme = tmp[sub.Type].(db.MemePost)
				} else {
					rMeme, err := getRandomMeme()
					if err != nil {
						return err
					}
					tmp[sub.Type] = rMeme
					meme = rMeme
				}
				err = sendSingleMeme(sub.ChatId, meme)
				break
			}

			sub.LastAction = nowTs
			err = db.SaveSchedule(sub)
		}
	}
	return err
}

func getFeedsList() string {
	var res []string
	for key, val := range feedsDescription {
		res = append(res, key+" - "+val)
	}
	return strings.Join(res, "\n")
}

func subscribe(msg telegram.Message) error {
	var err error
	regex := *regexp.MustCompile(`\s+(boobs|memes)(\s+(\d+))?`)
	match := regex.FindStringSubmatch(msg.Text)

	title := msg.User.FirstName + " @" + msg.User.UserName
	isGroup := false

	if msg.Chat.Type != "private" {
		title = msg.Chat.Title
		isGroup = true
	}

	if len(match) == 4 {
		var period int
		if match[3] != "" {
			period, err = strconv.Atoi(match[3])
		} else {
			period = MinPeriod
		}

		if period < MinPeriod {
			period = MinPeriod
		}

		now := time.Now()
		nowTs := now.Unix()

		schedule := db.Schedule{
			Type:       match[1],
			Period:     period,
			ChatId:     msg.Chat.Id,
			Title:      title,
			LastAction: nowTs,
			IsGroup:    isGroup,
		}
		return db.SaveSchedule(schedule)
	} else {
		err = errors.New("Command not matched to pattern `/subscribe <feed> <period min>`")
	}
	return err
}

func unsubscribe(msg telegram.Message) error {
	var err error
	regex := *regexp.MustCompile(`\s+(boobs|memes)`)
	match := regex.FindStringSubmatch(msg.Text)

	if len(match) >= 2 {
		return db.RemoveFromSchedule(msg.Chat.Id, match[1])
	} else {
		err = errors.New("Command not matched to pattern `/subscribe <feed> <period min>`")
	}
	return err
}
