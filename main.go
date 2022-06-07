package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/smtp"
	"strings"
	"time"
)

var (
	sendFrom string
	sendTo   string
	password string
)

func loadConfig() {
	type A struct {
		From string `json:"send_from"`
		To   string `json:"send_to"`
		Pass string `json:"password"`
	}
	c, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	var res A
	err = json.Unmarshal(c, &res)
	if err != nil {
		log.Fatal(err)
	}
	sendFrom = res.From
	sendTo = res.To
	password = res.Pass
}

func sendEmailBySMTP(message string) error {
	body := strings.Join([]string{
		"To: " + sendTo,
		"Subject: テストメッセージ",
		message,
		"送信時刻: " + time.Now().String(),
	}, "\r\n")
	auth := smtp.PlainAuth(
		"",
		sendFrom,
		password,
		"smtp.gmail.com",
	)
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth, sendFrom,
		[]string{sendTo},
		[]byte(body),
	)
	if err != nil {
		log.Println(err)
		return errors.New("cannot send email")
	}
	return nil
}

func sendEmail(message string) error {
	return sendEmailBySMTP(message)
}

func main() {
	loadConfig()
	sendEmail("テスト送信です。")
}
