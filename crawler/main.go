package main

import (
	"log"
)

const authToken = "remember_82e5d2c56bdd0811318f0cf078b78bfc=eyJpdiI6InZCSFNvOWJnOE1NUDJQRzNnbXhqdlE9PSIsInZhbHVlIjoiNzVEZnhpZHFlRXhDQ1wvNzFJZkNkQlVjM0xTMmUydG84OEd6THRJQkc3cWxvOHhqZHhoTU1RNjg5MHBTdkhUa1dvelhnVEI4eVBielp3VzkxRDhBcXRYUDJlWWU2RHhMaW50eklwVHpjSGc4PSIsIm1hYyI6ImI1OTYyNjUzYjA0OTNlZDZjNjY5M2NhNDcyMzkwYTEwMmZlODQ5OGJjMGYwZTUyMTM2M2QzZWQyYzc1ZjdlYzkifQ%3D%3D; quiz_605=1; current_block=18; XSRF-TOKEN=eyJpdiI6IjM3cVhuR1VkWmVVV214RGZJOVo0Q2c9PSIsInZhbHVlIjoiTGowcFg2cnRVblhtNkVlRVZUN2hNODhORXViQ09HZkVPMnVtalwvVmVmRGJwXC9tdGVwVXlZMG1zTnB5UklOSzBJakNsQ0VmNm9laGc1aHJ1ZmQxcmgzQT09IiwibWFjIjoiMGEzM2E3NGI5ZDY2ZGJkNTJjMWMxODE1NTg4YzI2YzI0MTRjNzY4NTUyOTdjYWNlODFjYWFkZGYwNWMxOTJjYyJ9;"

var themeIDs = []int{130}

func main() {
	client := newClient(authToken)

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