package main

import (
	"github.com/pkg/errors"
	"main/db"
	"main/telegram"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func schedule() {
	db.SelectSchedule()
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

	if msg.Chat.Type != "private" {
		title = msg.Chat.Title
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
		}
		return db.SaveSchedule(schedule)
	} else {
		err = errors.New("Command not matched to pattern `/subscribe <feed> <period min>`")
	}
	return err
}
