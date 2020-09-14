package main

import (
	"flag"
	"log"
	"os"
)

var (
	themeIDs = []int{130}

	workDir string
)

func init() {
	flag.StringVar(&workDir, "dir", "", "directory for storing data (temporary directory used)")
}

func main() {
	flag.Parse()
	log.Printf("Use token: %s", authToken)
	log.Printf("Use work dir: %s", workDir)

	client := newClient()

	log.Println("Start crawling themes and questions")
	for _, id := range themeIDs {
		theme, err := client.FetchTheme(id)
		if err != nil {
			log.Printf("Error: %s", err)
			continue
		}
		log.Printf("Response theme: %+v", theme)

		questions, err := client.FetchQuestions(id)
		if err != nil {
			log.Printf("Error: %s", err)
			continue
		}
		log.Printf("Response questions: %+v", questions)
	}
}
