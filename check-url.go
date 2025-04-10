package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
)

const checkMark = "\u2713"

func main() {
	duration, frequency, url := processCommandLineArguments()
	log.Println("== start ==")
	log.Println("URL: ", url, checkMark)
	log.Println("Duration: ", duration, "seconds", checkMark)
	log.Println("Frequency: every", frequency, checkMark)
	color.Green("Hello")
	trigger := time.NewTicker(time.Duration(frequency) * time.Second)

	counter := 0
	ok := 0
	ko := 0
	go func() {
		for range trigger.C {
			if sendHttpRequest(url) {
				ok++
			} else {
				ko++
			}
			counter++
		}
	}()

	time.Sleep(time.Duration(duration) * time.Second)
	trigger.Stop()

	log.Println(ok, "/", counter, "requests OK")
	log.Println(ko, "/", counter, "requests NOT OK")
}

func sendHttpRequest(url string) (requestOk bool) {
	response, err := http.Get(url)
	if response == nil {
		log.Println(err)
		return false
	} else {
		log.Println("URL: ", url, "Response status:", response.StatusCode)
		return response.StatusCode == 200
	}
}

func processCommandLineArguments() (duration int, frequency int, url string) {
	var help bool

	flag.IntVar(&duration, "duration", 30, "Duration (in seconds) how long should the program run")
	flag.IntVar(&frequency, "frequency", 3, "How often should the program send the request, defaults to 'once every 3 seconds'")
	flag.StringVar(&url, "url", "http://localhost:8080/api/mytest", "URL where the requests should be sent")
	flag.BoolVar(&help, "help", false, "Program help")

	flag.Parse()

	if help {
		flag.Usage()
	}
	return duration, frequency, url
}
