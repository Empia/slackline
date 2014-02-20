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
	json, err := json.Marshal(s)
	if err != nil {
		return
	}
	//res, err := http.PostForm(postMessageURL, url.Values{"payload": {json}})
	content := []byte("payload=")
	content = append(content, json...)
	reader := bytes.NewReader(content)
	res, err := http.Post(
		postMessageURL+os.Getenv("SLACK_TOKEN"),
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
	m.Get("/", func() string {
		message := slackMessage{"#echo", "ernesto", "ernesto, probando, un dos tres"}
		message.send()
		msg, err := message.json()
		if err != nil {
			return err.Error()
		} else {
			return msg
		}
		//return "Hello world!"
	})
	m.Post("/post", func(res http.ResponseWriter, req *http.Request) {
		msg := slackMessage{
			Channel:  "#echo",
			Username: req.PostFormValue("user_name"),
			Text:     req.PostFormValue("text"),
		}
		fmt.Printf("Received %#v\n", msg)
		err := msg.send()
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			res.WriteHeader(500)
		} else {
			fmt.Println("Sent")
		}
	})
	m.Run()
}
