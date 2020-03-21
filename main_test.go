package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func removeWhitespace(s string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(s, " ", ""),
			"\t", ""),
		"\n", "")
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestItUpdatesTheDNSARecords(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method == http.MethodGet {
			return &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
	{
		"data": "100.100.100.100",
		"name": "*",
		"ttl": 600,
		"type": "A"
	},
	{
		"data": "100.100.100.100",
		"name": "@",
		"ttl": 600,
		"type": "A"
	}
]`)),
				Header: make(http.Header),
			}
		} else if req.Method == http.MethodPut {
			assert.Equal(t, req.URL.String(), "localhost:3333/v1/domains/example.com/records/A")
			assert.Equal(t, req.Method, http.MethodPut)
			b, err := ioutil.ReadAll(req.Body)
			assert.NoError(t, err)
			assert.Equal(t, string(b), removeWhitespace(`[
	{
		"data": "101.101.101.101",
		"name": "*",
		"ttl": 600,
		"type": "A"
	},
	{
		"data": "101.101.101.101",
		"name": "@",
		"ttl": 600,
		"type": "A"
	}
]`))
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		} else {
			return &http.Response{
				StatusCode: 500,
				Header:     make(http.Header),
			}
		}
	})

	runner := &Updater{}
	_, err := runner.CheckAndUpdate("example.com", "101.101.101.101", WithEndpoint("localhost:3333"), WithHttpClient{Client: client})
	if err != nil {
		t.Fatal(err)
	}
}

func TestItUpdatesAllFoundNames(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method == http.MethodGet {
			return &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(removeWhitespace(`[
	{
		"data": "100.100.100.100",
		"name": "*",
		"ttl": 600,
		"type": "A"
	},
	{
		"data": "100.100.100.100",
		"name": "@",
		"ttl": 600,
		"type": "A"
	},
	{
		"data": "100.100.100.100",
		"name": "bob",
		"ttl": 600,
		"type": "A"
	}
]`))),
				Header: make(http.Header),
			}
		} else if req.Method == http.MethodPut {
			assert.Equal(t, req.URL.String(), "localhost:3333/v1/domains/example.com/records/A")
			assert.Equal(t, req.Method, http.MethodPut)
			b, err := ioutil.ReadAll(req.Body)
			assert.NoError(t, err)
			assert.Equal(t, string(b), removeWhitespace(`[
	{
		"data": "101.101.101.101",
		"name": "*",
		"ttl": 600,
		"type": "A"
	},
	{
		"data": "101.101.101.101",
		"name": "@",
		"ttl": 600,
		"type": "A"
	},
	{
		"data": "101.101.101.101",
		"name": "bob",
		"ttl": 600,
		"type": "A"
	}
]`))

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		} else {
			return &http.Response{
				StatusCode: 500,
				Header:     make(http.Header),
			}
		}
	})

	runner := &Updater{}
	_, err := runner.CheckAndUpdate("example.com", "101.101.101.101", WithEndpoint("localhost:3333"), WithHttpClient{Client: client})
	if err != nil {
		t.Fatal(err)
	}
}

func TestItReportsSuccessWhenItUpdates(t *testing.T) {
	t.Skip("not implemented")
}

func TestItReportsFAilureWhenItFailsToUpdate(t *testing.T) {
	t.Skip("not implemented")
}

func TestItReportsNoChangeWhenNoChangeIsNeeded(t *testing.T) {
	t.Skip("not implemented")
}

func TestItOnlyUpdatesIfItDetectsAChange(t *testing.T) {
	t.Skip("not implemented")
}
