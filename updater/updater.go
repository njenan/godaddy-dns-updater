package updater

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
	RecordNames map[string]bool
	DryRun      bool
	AuthKey     string
	AuthSecret  string
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
	c.RecordNames[string(d)] = true
	return c
}

var _ Option = WithRecordName("")

type WithDryRun bool

func (d WithDryRun) Apply(c Config) Config {
	c.DryRun = bool(d)
	return c
}

var _ Option = WithDryRun(true)

type WithAuthKey string

func (k WithAuthKey) Apply(c Config) Config {
	c.AuthKey = string(k)
	return c
}

var _ Option = WithAuthKey("")

type WithAuthSecret string

func (s WithAuthSecret) Apply(c Config) Config {
	c.AuthSecret = string(s)
	return c
}

var _ Option = WithAuthSecret("")

type Report struct {
	DidUpdate bool
	Records   []*Record
}

type Record struct {
	Data string `json:"data"`
	Name string `json:"name"`
	TTL  int    `json:"ttl"`
	Type string `json:"type"`
}

type Updater struct {
}

func (r *Updater) CheckAndUpdate(domain, targetIP string, options ...Option) (*Report, error) {
	// Set up config default
	c := Config{
		Endpoint:    "https://api.godaddy.com",
		HttpClient:  http.DefaultClient,
		RecordNames: make(map[string]bool),
	}

	for _, o := range options {
		c = o.Apply(c)
	}

	url := c.Endpoint + "/v1/domains/" + domain + "/records/A"
	req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	if c.AuthKey != "" {
		req.Header.Add("Authorization", "sso-key "+c.AuthKey+":"+c.AuthSecret)
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

	var updateNeeded bool
	for _, r := range records {
		if len(c.RecordNames) == 0 || c.RecordNames[r.Name] {
			if r.Data != targetIP {
				updateNeeded = true
				r.Data = targetIP
			}
		}
	}

	report := &Report{Records: records}
	if !c.DryRun && updateNeeded {
		b, err := json.Marshal(records)
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest(http.MethodPut, url, strings.NewReader(string(b)))
		if err != nil {
			return nil, err
		}

		if c.AuthKey != "" {
			req.Header.Add("Authorization", c.AuthKey+":"+c.AuthSecret)
		}

		_, err = c.HttpClient.Do(req)
		if err != nil {
			return nil, err
		}

		report.DidUpdate = true
	}

	return report, err
}
