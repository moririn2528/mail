package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var (
	sendFrom string
	sendTo   string
	password string
)

const region = "us-west-2"

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

func sendEmailBySES(message string) error {
	sendErr := errors.New("cannot send email")

	from := "testmail@e.torimari.site"
	to := "strangenoise1@gmail.com"
	subject := "テスト"
	data := strings.Join([]string{
		message,
		"送信時刻: " + time.Now().String(),
	}, "\r\n")

	awsAccessKey := "AKIAQBZGMKAFAZCT6GOT"
	awsSecretKey := "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"
	awsSession := session.New(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	})
	client := ses.New(awsSession)
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(data),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(from),
	}
	_, err := client.SendEmail(input)
	if err != nil {
		log.Println(err)
		return sendErr
	}
	return nil
}

func sendEmail(message string) error {
	_, local := os.LookupEnv("windir")
	if local {
		return sendEmailBySMTP(message)
	} else {
		return sendEmailBySES(message)
	}
}

func main() {
	loadConfig()
	sendEmail("テスト送信です。")
}
