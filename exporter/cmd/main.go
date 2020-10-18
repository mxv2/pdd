package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"exporter"
)

func main() {
	const database = "test.db"
	//const database = ":memory:"
	db := sqlx.MustOpen("sqlite3", database)

	ctx := context.Background()
	_, err := sqlx.LoadFileContext(ctx, db, "init.sql")
	if err != nil {
		panic(err)
	}

	log.Println("Init db: success")

	const dataDir = "../testdata"
	const modelID = 1592507431253
	const deckID = 1595179860403

	stmt, err := db.PrepareNamedContext(ctx, `INSERT INTO notes (id, guid, mid, mod, usn, tags, flds, sfld, csum, flags, data) 
VALUES (:mod, :guid, :mid, :mod, -1, '', :flds, :sfld, :csum, 0, '')`)
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

		image := filepath.Join(path, "image.jpg")
		if _, err := os.Stat(filepath.Join(path, "image.jpg")); os.IsNotExist(err) {
			image = ""
		}

		answer := ""
		noteWriter := &exporter.NoteWriter{}

		noteWriter.WriteImage(image)

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

		type dto struct {
			ID            int64  `db:"mod"`
			GUID          string `db:"guid"`
			ModelID       int64  `db:"mid"`
			Fields        string `db:"flds"`
			StripedFields string `db:"sfld"`
			Checksum      uint32 `db:"csum"`
		}
		_, err = stmt.ExecContext(ctx, dto{
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
		return filepath.SkipDir
	})
	if err != nil {
		log.Fatalf("Add notes: failed cause %s", err)
	}
	log.Println("Add notes: success")
}
