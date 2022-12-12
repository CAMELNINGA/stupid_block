package main

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
)

type Client struct {
	baseUrl string
	client  *resty.Client
}

func NewClient(baseUrl string) Client {
	client := resty.New()
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	client.JSONMarshal = json.Marshal
	client.JSONUnmarshal = json.Unmarshal
	c := Client{
		baseUrl: baseUrl,
		client:  client,
	}
	return c
}

type Sign struct {
	Status         int    `json:"status"`
	StatusString   string `json:"statusString"`
	TimeStampToken struct {
		Ts        string `json:"ts"`
		Signature string `json:"signature"`
	} `json:"timeStampToken"`
}

func (c Client) GetSign(hash string) (string, time.Time) {
	u, err := url.Parse(c.baseUrl + "/ts")
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	q.Set("digest", hash)
	u.RawQuery = q.Encode()
	resp, err := c.client.R().EnableTrace().Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	var sign Sign
	if err := json.Unmarshal(resp.Body(), &sign); err != nil {
		log.Fatal("Client unmarshal failed: " + err.Error())
	}
	date, err := time.Parse("2006-01-02T15:04:05.000-07", sign.TimeStampToken.Ts)
	if err != nil {
		log.Fatal("Client unmarshal failed: " + err.Error())
	}
	return sign.TimeStampToken.Signature, date
}
