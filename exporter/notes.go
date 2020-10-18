package exporter

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"io"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

const separator = "\u001F"

type Note struct {
	id     int64
	deckID int64
	front  string
	back   string
}

func NewNote(deckID int64, front string, back string) *Note {
	return &Note{
		id:     int64(time.Now().Nanosecond()),
		deckID: deckID,
		front:  front,
		back:   back,
	}
}

func (n Note) ID() int64 {
	return n.id
}

func (n Note) GUID() string {
	const table = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&()*+,-./:;<=>?@[]^_`{|}~"
	const rank = int64(len(table))
	r := rand.NewSource(time.Now().UnixNano())
	i := r.Int63()
	var buf bytes.Buffer
	for i > 0 {
		q, r := i/rank, i%rank
		buf.WriteByte(table[r])
		i = q
	}
	for buf.Len() < 10 {
		buf.WriteByte(0)
	}
	return buf.String()
}

func (n Note) Checksum() uint32 {
	h := sha1.New()
	_, _ = io.WriteString(h, stripHtmlPreservingImage(n.front))
	digest := h.Sum(nil)
	return binary.BigEndian.Uint32(digest[:4])
}

func (n Note) StripedFront() string {
	return stripHtmlPreservingImage(n.front)
}

func (n Note) Fields() string {
	return n.front + separator + n.back
}

var htmlEntitiesRe = regexp.MustCompile("<.*?>")
var srcAttributeRe = regexp.MustCompile("src=\"(.*?)\"")

func stripHtmlPreservingImage(s string) string {
	return htmlEntitiesRe.ReplaceAllStringFunc(s, func(occurrence string) string {
		if strings.HasPrefix(occurrence, "<img") {
			match := srcAttributeRe.FindAllStringSubmatch(occurrence, -1)
			if len(match) > 0 && len(match[0]) > 1 {
				return " " + match[0][1] + " "
			}
		}
		return ""
	})
}

type NoteWriter struct {
	buf bytes.Buffer
}

func (w *NoteWriter) WriteString(s string) {
	w.buf.WriteString("<div>" + s + "</div>")
}

func (w *NoteWriter) WriteNewLine() {
	w.buf.WriteString("<div><br></div>")
}

func (w *NoteWriter) WriteImage(image string) {
	w.buf.WriteString("<img src=\"" + image + "\">")
}

func (w *NoteWriter) String() string {
	return w.buf.String()
}

func (w *NoteWriter) Reset() {
	w.buf.Reset()
}
