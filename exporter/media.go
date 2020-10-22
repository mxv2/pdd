package exporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"
)

type AddedImage struct {
	hash     string
	realpath string
}

type Storage struct {
	added []AddedImage
}

func (s *Storage) AddImage(path string) string {
	hash := guid() + ".jpg"
	s.added = append(s.added, AddedImage{hash: hash, realpath: path})
	return hash
}

func (s *Storage) WriteFiles(dir string) error {
	for i, img := range s.added {
		source, err := os.Open(img.realpath)
		if err != nil {
			return fmt.Errorf("media: failed open image file %s: %w", img.realpath, err)
		}

		filename := strconv.FormatInt(int64(i), 10)
		filepath := path.Join(dir, filename)
		target, err := os.Create(filepath)
		if err != nil {
			return fmt.Errorf("media: failed create image file %s: %w", filepath, err)
		}

		_, err = io.Copy(target, source)
		if err != nil {
			return fmt.Errorf("media: failed copy image data from %s: %w", source.Name(), err)
		}

		target.Close()
		source.Close()
	}
	return nil
}

func (s *Storage) WriteHashFile(dir string) error {
	hash := make(map[string]string, len(s.added))
	for i, img := range s.added {
		index := strconv.FormatInt(int64(i), 10)
		hash[index] = img.hash
	}

	data, err := json.Marshal(hash)
	if err != nil {
		return fmt.Errorf("media: failed marshal hash file: %w", err)
	}

	const filename = "media"
	filepath := path.Join(dir, filename)

	err = ioutil.WriteFile(filepath, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("media: failed write hash file: %w", err)
	}

	return nil
}

func guid() string {
	const table = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const rank = int64(len(table))
	r := rand.NewSource(time.Now().UnixNano())
	i := r.Int63()
	var buf bytes.Buffer
	for i > 0 {
		q, r := i/rank, i%rank
		buf.WriteByte(table[r])
		i = q
	}
	for buf.Len() < 8 {
		buf.WriteByte(0)
	}
	return buf.String()
}
