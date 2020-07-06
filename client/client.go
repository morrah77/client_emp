package client

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

const MAX_ATTEMPTS  = 3

type Client struct {
	listUrl string
	itemUrl string
	maxAttempts int
	itemIdValidator func(s string) bool
}

// TODO add urls validation, return error if invalid
func NewClient(listUrl string, itemUrl string, maxAttempts int, itemValidator func(s string) bool) (*Client, error)  {
	if maxAttempts == 0 {
		maxAttempts = MAX_ATTEMPTS
	}
	if itemValidator == nil {
		itemValidator = defaultItemValidator
	}
	return &Client{
		listUrl:         listUrl,
		itemUrl:         itemUrl,
		maxAttempts:    maxAttempts,
		itemIdValidator: itemValidator,
	}, nil
}

func(c *Client) FetchList() ([]byte, error) {
	return c.fetchResult(c.listUrl)
}

func(c *Client) FetchItem(id string) ([]byte, error) {
	if !(c.itemIdValidator(id)) {
		return nil, errors.New("Invalid ID")
	}
	url := fmt.Sprintf(c.itemUrl, id)
	return c.fetchResult(url)
}

func(c *Client) fetchResult(url string) (res []byte, err error) {
	var (
		status int
		resp *http.Response
	)
	for i:= 0; (i < c.maxAttempts) && (status != http.StatusOK); i++ {
		resp, err = http.Get(url)
		if err != nil {
			//fmt.Printf("An error occured: %s\n", err.Error())
			return nil, err
		}
		status = resp.StatusCode
	}
	p := &bytes.Buffer{}
	_, err = p.ReadFrom(resp.Body)
	if (err == nil) && status != http.StatusOK {
		err = errors.New(resp.Status)
	}
	return p.Bytes(), err
}

func defaultItemValidator(s string) bool {
	return true
}