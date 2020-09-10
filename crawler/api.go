package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36"

type client struct {
	host      string
	authToken string
	tClient   *http.Client
}

func newClient(authToken string) *client {
	return &client{
		host:      "https://profteh.com",
		authToken: authToken,
		tClient: &http.Client{
			CheckRedirect: func(*http.Request, []*http.Request) error { return fmt.Errorf("no redirects") },
		},
	}
}

type theme struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Tag   string `json:"slug"`
}

func (c *client) FetchTheme(id int) (theme, error) {
	url := fmt.Sprintf("%s/%s/%d", c.host, "api/get/type", id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return theme{}, err
	}
	req.Header.Add("Cookie", c.authToken)
	req.Header.Add("User-Agent", userAgent)

	resp, err := c.tClient.Do(req)
	if err != nil {
		return theme{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return theme{}, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return theme{}, err
	}

	var target theme
	err = json.Unmarshal(data, &target)
	return target, nil
}

type question struct {
	ID      int     `json:"id"`
	Title   string  `json:"title"`
	Options options `json:"answers"`
	Answer  int     `json:"answerNum"`
	Image   string  `json:"image_alt_imgur"`
}

type options []string

func (o *options) UnmarshalJSON(data []byte) error {
	unquote, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	var t []string
	err = json.Unmarshal([]byte(unquote), &t)
	if err != nil {
		return err
	}
	*o = t
	return nil
}

func (c *client) FetchQuestions(id int) ([]question, error) {
	url := fmt.Sprintf("%s/%s/%d", c.host, "api/get/exam", id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []question{}, err
	}
	req.Header.Add("Cookie", c.authToken)
	req.Header.Add("User-Agent", userAgent)

	resp, err := c.tClient.Do(req)
	if err != nil {
		return []question{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return []question{}, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []question{}, err
	}

	var target []question
	err = json.Unmarshal(data, &target)
	if err != nil {
		return []question{}, err
	}
	return target, nil
}