package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"exporter"
)

const (
	database = "collection.anki2"
)

var (
	data = flag.String("data", "data", "source data directory")
	out  = flag.String("out", "out", "directory for output data")
	name = flag.String("name", "PDD.apkg", "name of the target .apkg file")

	model = flag.Int64("model", 1592507431253, "model - target model ID in Anki")
	deck  = flag.Int64("deck", 1595179860403, "deck - target deck ID in Anki")
)

func main() {
	flag.Parse()
	if out == nil || *out == "" {
		log.Fatalf("out should not be empty")
	}
	outputDir := *out
	if data == nil || *data == "" {
		log.Fatalf("data should not be empty")
	}
	dataDir := *data
	if name == nil || *name == "" {
		log.Fatalf("name should not be empty")
	}
	outputName := *name
	if model == nil || *model == 0 {
		log.Fatalf("model should not be 0")
	}
	modelID := *model
	if deck == nil || *deck == 0 {
		log.Fatalf("deck should not be 0")
	}
	deckID := *deck

	dbFile := path.Join(outputDir, database)
	db := sqlx.MustOpen("sqlite3", dbFile)
	_, err := sqlx.LoadFile(db, "init.sql")
	if err != nil {
		log.Fatalf("Init db: %s", err)
	}
	noteStmt, err := db.PrepareNamed(`INSERT INTO notes (id, guid, mid, mod, usn, tags, flds, sfld, csum, flags, data) 
VALUES (:mod, :guid, :mid, :mod, -1, '', :flds, :sfld, :csum, 0, '')`)
	if err != nil {
		log.Fatalf("Init db: %s", err)
	}
	cardStmt, err := db.PrepareNamed(`INSERT INTO cards(nid, did, ord, mod, usn, type, queue, due, ivl, factor, reps, lapses, "left", odue, odid, flags, data)
VALUES (:nid, :did, 0, :mod, -1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, '');`)
	if err != nil {
		log.Fatalf("Init db: %s", err)
	}
	log.Println("Init db: success")

	storage := &exporter.Storage{}
	err = filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() || !strings.HasPrefix(info.Name(), "q") {
			return nil
		}

		file, err := os.Open(filepath.Join(path, "question.txt"))
		if err != nil {
			return err
		}
		defer file.Close()

		imagePath := filepath.Join(path, "image.jpg")
		_, err = os.Stat(imagePath)
		if err != nil {
			imagePath = ""
		}

		answer := ""
		noteWriter := &exporter.NoteWriter{}

		var imageID string
		if imagePath != "" {
			imageID = storage.AddImage(imagePath)
		}
		if imageID != "" {
			noteWriter.WriteImage(imageID)
			noteWriter.WriteNewLine()
		}

		var block string
		r := bufio.NewScanner(file)
		for r.Scan() {
			line := r.Text()

			if block == "A:" && line == "" {
				break
			}

			if block == "" || line == "A:" {
				block = line
				continue
			}

			switch block {
			case "Q:":
				if line == "" {
					noteWriter.WriteNewLine()
				} else {
					noteWriter.WriteString(line)
				}
			case "A:":
				answer += line
			}
		}

		note := exporter.NewNote(deckID, noteWriter.String(), answer)

		type noteDTO struct {
			ID            int64  `db:"mod"`
			GUID          string `db:"guid"`
			ModelID       int64  `db:"mid"`
			Fields        string `db:"flds"`
			StripedFields string `db:"sfld"`
			Checksum      uint32 `db:"csum"`
		}
		_, err = noteStmt.Exec(noteDTO{
			ID:            note.ID(),
			GUID:          note.GUID(),
			ModelID:       modelID,
			Fields:        note.Fields(),
			StripedFields: note.StripedFront(),
			Checksum:      note.Checksum(),
		})
		if err != nil {
			return err
		}

		type cardDTO struct {
			ID     int64 `db:"mod"`
			DeckID int64 `db:"did"`
			NoteID int64 `db:"nid"`
		}
		_, err = cardStmt.Exec(cardDTO{
			ID:     time.Now().Unix(),
			DeckID: deckID,
			NoteID: note.ID(),
		})
		if err != nil {
			return err
		}

		return filepath.SkipDir
	})
	if err != nil {
		log.Fatalf("Add notes: failed cause %s", err)
	}
	log.Println("Add notes: success")

	err = storage.WriteFiles(outputDir)
	if err != nil {
		log.Fatalf("Add images: failed cause %s", err)
	}
	log.Println("Add images: success")

	err = storage.WriteHashFile(outputDir)
	if err != nil {
		log.Fatalf("Add images index: failed cause %s", err)
	}
	log.Println("Add images index: success")

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	err = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		source, err := os.Open(path)
		if err != nil {
			return err
		}
		defer source.Close()

		target, err := w.Create(info.Name())
		if err != nil {
			return err
		}

		_, err = io.Copy(target, source)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(path.Join(outputDir, outputName), buf.Bytes(), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Write zip archive: success")
}
