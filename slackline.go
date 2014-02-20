package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codegangsta/martini"
	"io"
	"io/ioutil"
	"net/http"
)

const postMessageURL = "/services/hooks/incoming-webhook?token="

type slackMessage struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

func (s slackMessage) payload() io.Reader {
	content := []byte("payload=")
	json, _ := json.Marshal(s)
	content = append(content, json...)
	return bytes.NewReader(content)
}

func (s slackMessage) sendTo(domain, token string) (err error) {
	payload := s.payload()

	res, err := http.Post(
		"https://"+domain+postMessageURL+token,
		"application/x-www-form-urlencoded",
		payload,
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
	m.Post("/bridge", func(res http.ResponseWriter, req *http.Request) {
		username := req.PostFormValue("user_name")
		text := req.PostFormValue("text")

		if username == "slackbot" {
			// Avoid infinite loop
			return
		}

		msg := slackMessage{
			Username: username,
			Text:     text,
		}

		domain := req.URL.Query().Get("domain")
		token := req.URL.Query().Get("token")

		err := msg.sendTo(domain, token)

		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			res.WriteHeader(500)
		} else {
			fmt.Println("Message sent")
		}
	})
	m.Run()
}
