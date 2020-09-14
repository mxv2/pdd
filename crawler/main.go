package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
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
	if workDir == "" {
		tmpDir, err := ioutil.TempDir("", "crawler")
		if err != nil {
			log.Fatal("Can not obtain tmp directory")
		}
		workDir = tmpDir
	}
	log.Printf("Use work dir: %s", workDir)

	client := newClient()

	for _, id := range themeIDs {
		log.Printf("Start crawling theme and questions for ID: %d", id)
		theme, err := client.FetchTheme(id)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		log.Printf("Response theme: %+v", theme)

		curDir := path.Join(workDir, theme.Tag)
		err = os.MkdirAll(curDir, 0777)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}

		log.Printf("Create info.txt file")
		infoFile, err := os.Create(path.Join(curDir, "info.txt"))
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		log.Printf("Write to info.txt file")
		_, err = infoFile.WriteString(theme.Title + "\n")
		if err != nil {
			log.Fatalf("Error: %s", err)
		}

		questions, err := client.FetchQuestions(id)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		log.Printf("Response questions: %+v", questions)

		for i, q := range questions {
			qDir := path.Join(curDir, fmt.Sprintf("q%d", i+1))
			err = os.MkdirAll(qDir, 0777)
			if err != nil {
				log.Fatalf("Error: %s", err)
			}

			qFile, err := os.Create(path.Join(qDir, "question.txt"))
			if err != nil {
				log.Fatalf("Error: %s", err)
			}

			_, err = qFile.WriteString("Q:\n")
			if err != nil {
				log.Fatalf("Error: %s", err)
			}
			_, err = qFile.WriteString(q.Title + "\n\n")
			if err != nil {
				log.Fatalf("Error: %s", err)
			}
			for _, opt := range q.Options {
				_, err = qFile.WriteString("* " + opt + "\n")
				if err != nil {
					log.Fatalf("Error: %s", err)
				}
			}
			_, err = qFile.WriteString("\nA:\n")
			if err != nil {
				log.Fatalf("Error: %s", err)
			}
			_, err = qFile.WriteString(q.Options[q.Answer] + "\n")
			if err != nil {
				log.Fatalf("Error: %s", err)
			}

			if q.Image != "" {
				err := writeImage(q.Image, path.Join(qDir, "image.jpg"))
				if err != nil {
					log.Fatalf("Error: %s", err)
				}
			}
		}
	}
}

var imageLoader = http.Client{}

func writeImage(imageURL string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := imageLoader.Get(imageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
