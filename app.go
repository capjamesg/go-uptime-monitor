package main

import (
	"os"
	"net/http"
	"net/smtp"
	"time"
	"log"
	"github.com/joho/godotenv"
	"fmt"
	"sync"
)

type Email struct {
	Sender string
	Password string
	Host string
	Port string
	SendTo []string
}

func sendEmail (url string) {
	email := Email{
		os.Getenv("GMAIL_USERNAME"),
		os.Getenv("GMAIL_PASSWORD"),
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		[]string{os.Getenv("SEND_TO")},
	}

	fmt.Println("Sending email to " + email.SendTo[0] + "..." + url)

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s service is down\r\n\r\n%s is down as of %s", email.Sender, email.SendTo[0], url, 
url, time.Now().String())

	smtpData := []byte(message)

	authentication := smtp.PlainAuth("", email.Sender, email.Password, email.Host)

	err := smtp.SendMail(email.Host+":"+email.Port, authentication, email.Sender, email.SendTo, smtpData)

	if err != nil {
		log.Fatal(err)
		log.Fatal("Something went wrong sending an email.")
	}
}

func fetchData (service string) {
	response, err := http.Get(service)

	if err != nil {
		fmt.Println(service + " is down.")
		sendEmail(service)

		return
	} else {
		defer response.Body.Close()

		if response.StatusCode == 500 {
			sendEmail(service)
		}

		return
	}
}

func main () {
	err := godotenv.Load("/home/james/canary/.env")

	if err != nil {
		log.Fatal(err)
	}

	services := []string{
		"https://jamesg.blog",
		"https://indieweb-search.jamesg.blog",
		"https://es-indieweb-search.jamesg.blog",
		"https://breakfastand.coffee",
		"https://jamesg.blog/search/",
		"https://cali.jamesg.blog/",
		"https://etherpad.jamesg.blog",
		"https://grafana.jamesg.blog",
		"https://create.breakfastand.coffee",
		"https://vouch.breakfastand.coffee/auth",
		"https://coffeepot.jamesg.blog",
		"https://task.jamesg.blog",
		"https://sparkline.jamesg.blog",
		"https://jamesg.coffee",
	}

	var waitGroup sync.WaitGroup

	for _, service := range services {
		waitGroup.Add(1)

		fmt.Println("Checking " + service + " ...")

		go func(service string) {
			defer waitGroup.Done()
			fetchData(service)
		}(service)
	}

	waitGroup.Wait()
}
