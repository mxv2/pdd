package main

import (
	"bufio"
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

func main() {
	const outputDir = "output"
	_, err := os.Stat(outputDir)
	if err != nil {
		os.RemoveAll(outputDir)
	}
	if err != nil || os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	database := path.Join(outputDir, "collection.anki2")
	//const database = ":memory:"
	db := sqlx.MustOpen("sqlite3", database)

	_, err = sqlx.LoadFile(db, "init.sql")
	if err != nil {
		panic(err)
	}

	log.Println("Init db: success")

	const dataDir = "../testdata"
	const modelID = 1592507431253
	const deckID = 1595179860403

	storage := &exporter.Storage{}

	noteStmt, err := db.PrepareNamed(`INSERT INTO notes (id, guid, mid, mod, usn, tags, flds, sfld, csum, flags, data) 
VALUES (:mod, :guid, :mid, :mod, -1, '', :flds, :sfld, :csum, 0, '')`)
	if err != nil {
		panic(err)
	}

	cardStmt, err := db.PrepareNamed(`INSERT INTO cards(nid, did, ord, mod, usn, type, queue, due, ivl, factor, reps, lapses, "left", odue, odid, flags, data)
VALUES (:nid, :did, 0, :mod, -1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, '');`)
	if err != nil {
		panic(err)
	}

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

		imagePath, err := filepath.Abs(filepath.Join(path, "image.jpg"))
		if err != nil {
			return err
		}
		_, err = os.Stat(imagePath)
		if os.IsNotExist(err) {
			imagePath = ""
		}
		if err != nil {
			return err
		}

		answer := ""
		noteWriter := &exporter.NoteWriter{}

		var imageID string
		if imagePath != "" {
			imageID = storage.AddImage(imagePath)
		}
		if imageID != "" {
			noteWriter.WriteImage(imageID)
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
				noteWriter.WriteString(line)
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
}
