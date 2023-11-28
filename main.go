package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func do(visaProcessorUrl string) {
	var client http.Client = http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(visaProcessorUrl, "application/json", strings.NewReader("{}"))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)
	log.Printf("Response from visa processor: %s", resp.Status)
}

func main() {
	var visaProcessorUrl string
	var delay time.Duration = 60 * time.Second
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if pair[0] == "VISA_PROCESSOR_URL" {
			visaProcessorUrl = pair[1]
		}
		if pair[0] == "DELAY" {
			var err error
			delay, err = time.ParseDuration(pair[1])
			if err != nil {
				log.Fatalf("Cannot parse environment variable DELAY: %s", err)
			}
		}
	}

	if visaProcessorUrl == "" {
		log.Fatalln("No VISA_PROCESSOR_URL environment variable!")
	}

	log.Printf("Visa Processor is at URL '%s'", visaProcessorUrl)

	for {
		do(visaProcessorUrl)
		time.Sleep(delay)
	}
}
