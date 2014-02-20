package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codegangsta/martini"
	"io/ioutil"
	"net/http"
	"os"
)

const postMessageURL = "https://espresso.slack.com/services/hooks/incoming-webhook?token="

type slackMessage struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

func (s slackMessage) json() (msg string, err error) {
	m, err := json.Marshal(s)
	if err != nil {
		return
	}
	msg = string(m[:])
	return
}

func (s slackMessage) send() (err error) {
	return s.sendTo(os.Getenv("SLACK_TOKEN"))
}

func (s slackMessage) sendTo(token string) (err error) {
	json, err := json.Marshal(s)
	if err != nil {
		return
	}
	//res, err := http.PostForm(postMessageURL, url.Values{"payload": {json}})
	content := []byte("payload=")
	content = append(content, json...)
	reader := bytes.NewReader(content)
	res, err := http.Post(
		postMessageURL+token,
		"application/x-www-form-urlencoded",
		reader,
	)
	if res.StatusCode != 200 {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		return errors.New(res.Status + " - " + string(body))
	}
	return
}

func main() {
	m := martini.Classic()
	m.Get("/", func(res http.ResponseWriter, req *http.Request) string {
		message := slackMessage{"", "ernesto", "ernesto, probando, un dos tres"}
		token := req.URL.Query().Get("token")
		message.sendTo(token)
		msg, err := message.json()
		if err != nil {
			return err.Error()
		} else {
			return msg
		}
		//return "Hello world!"
	})
	m.Post("/bridge", func(res http.ResponseWriter, req *http.Request) {
		msg := slackMessage{
			Username: req.PostFormValue("user_name"),
			Text:     req.PostFormValue("text"),
		}
		token := req.URL.Query().Get("token")
		err := msg.sendTo(token)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			res.WriteHeader(500)
		} else {
			fmt.Println("Sent")
		}
	})
	m.Run()
}
