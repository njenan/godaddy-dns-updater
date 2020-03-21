package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Option interface {
	Apply(Config) Config
}

type Config struct {
	Endpoint    string
	HttpClient  *http.Client
	RecordNames []string
}

type WithEndpoint string

func (e WithEndpoint) Apply(c Config) Config {
	c.Endpoint = string(e)
	return c
}

var _ Option = WithEndpoint("")

type WithHttpClient struct {
	*http.Client
}

func (h WithHttpClient) Apply(c Config) Config {
	c.HttpClient = h.Client
	return c
}

var _ Option = WithHttpClient{}

type WithRecordName string

func (d WithRecordName) Apply(c Config) Config {
	c.RecordNames = append(c.RecordNames, string(d))
	return c
}

type Record struct {
	Data string `json:"data"`
	Name string `json:"name"`
	TTL  int    `json:"ttl"`
	Type string `json:"type"`
}

type Updater struct {
}

func (r *Updater) CheckAndUpdate(domain, targetIP string, options ...Option) (interface{}, error) {
	c := Config{}
	for _, o := range options {
		c = o.Apply(c)
	}

	if c.HttpClient == nil {
		c.HttpClient = http.DefaultClient
	}

	url := c.Endpoint + "/v1/domains/" + domain + "/records/A"
	req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("got status code %v", resp.StatusCode)
	}

	var records []*Record
	err = json.NewDecoder(resp.Body).Decode(&records)
	if err != nil {
		return nil, err
	}

	for _, r := range records {
		r.Data = targetIP
	}

	b, err := json.Marshal(records)
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequest(http.MethodPut, url, strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}

	_, err = c.HttpClient.Do(req)

	return nil, err
}

func main() {
}
