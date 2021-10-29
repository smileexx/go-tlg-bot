package telegram

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var commands = map[string]string{
	"help":  "display all commands",
	"boobs": "send random ero image",
	"feeds": "list of available feeds to subscribe",
}

var help = map[string]string{
	"/help":            "display this message =)",
	"/tag #<tag_word>": "send random image with tag",
	"/post <12345678>": "send post by post id",
}

const (
	PathSetWebhook            = "/setWebhook"
	PathSendMessage           = "/sendMessage"
	PathSendPhoto             = "/sendPhoto"
	PathSendVideo             = "/sendVideo"
	PathSendMediaGroup        = "/sendMediaGroup"
	PathSetMyCommands         = "/setMyCommands"
	PathGetChatAdministrators = "/getChatAdministrators"
)

const API_URL = "https://api.telegram.org/bot"

func BuildUrl(param string) string {
	return API_URL + os.Getenv("BOT_TOKEN") + param
}

func HandleWebhook(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	body, _ := ioutil.ReadAll(req.Body)
	log.Println(string(body))
	var update Update
	err := json.Unmarshal(body, &update)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(update)
	err = handleMessage(update.Message)

	if err != nil {
		log.Fatal(err)
	}
}

func handleMessage(msg Message) error {
	return SendMessage(msg, "")
}

func SendPhoto(msg Message, srcUrl string, caption string) error {
	// outData := OutPhoto{ChatId: msg.Chat.Id, Photo: srcUrl}
	outPhoto := OutPhoto{
		ChatId:  msg.Chat.Id,
		Photo:   srcUrl,
		Caption: caption,
		// "entities":[{"offset":10,"length":4,"type":"hashtag"},{"offset":15,"length":48,"type":"url"}]}
	}
	// outPhoto.ReplayToMsgId = msg.Id
	return sendJson(PathSendPhoto, outPhoto)
}

func SendVideo(msg Message, srcUrl string, caption string) error {
	outVideo := OutVideo{
		ChatId:  msg.Chat.Id,
		Video:   srcUrl,
		Caption: caption,
	}
	// outVideo.ReplayToMsgId = msg.Id
	return sendJson(PathSendVideo, outVideo)
}

func SendMessage(msg Message, text string) error {
	if text == "" {
		text = ">>" + msg.Text
	}
	outData := OutMessage{
		ChatId:        msg.Chat.Id,
		Text:          text,
		ReplayToMsgId: msg.Id,
	}
	return sendJson(PathSendMessage, outData)
}

func SendMediaGroup(msg Message, media []InputMediaItem) error {
	outData := OutMediaGroup{
		ChatId: msg.Chat.Id,
		Media:  media,
	}
	return sendJson(PathSendMediaGroup, outData)
}

func SendHelp(msg Message) error {
	var res string
	for key, desc := range help {
		res += key + "\t - " + desc + "\n"
	}
	return SendMessage(msg, res)
}

func SetCommands() error {
	var commandList []BotCommand
	for key, desc := range commands {
		commandList = append(commandList, BotCommand{Command: key, Description: desc})
	}
	outData := OutMyCommands{
		Commands: commandList,
	}
	return sendJson(PathSetMyCommands, outData)
}

func sendJson(urlPath string, outData interface{}) error {
	body, err := json.Marshal(outData)
	if err != nil {
		log.Fatal(err)
		return err
	}
	bodyBytes := bytes.NewBuffer(body)
	log.Println(bodyBytes)
	_, err = http.Post(BuildUrl(urlPath), "application/json", bodyBytes)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func DeleteWebhook() {
	http.Get(BuildUrl(PathSetWebhook + "?url="))
}

func SetWebhook() {
	http.Get(BuildUrl(PathSetWebhook + "?url=" + os.Getenv("HOST") + os.Getenv("BOT_TOKEN")))
}

func getChatAdministrators(chatId int) ([]ChatMember, error) {
	resp, err := http.Get(BuildUrl(PathGetChatAdministrators + "?chat_id=" + strconv.Itoa(chatId)))
	var members ChatAdministrators
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &members)

	return members.Members, err
}

func IsUserAdmin(msg Message) (bool, error) {
	members, err := getChatAdministrators(msg.Chat.Id)
	if err != nil {
		return false, err
	}
	for _, user := range members {
		if user.User.Id == msg.User.Id {
			return true, err
		}
	}
	return false, err
}
