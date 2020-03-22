package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func recordStrings(r []Record) string {
	b, _ := json.Marshal(r)
	return string(b)
}

func TestItUpdatesTheDNSARecords(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method == http.MethodGet {
			return &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(recordStrings([]Record{
					{Data: "100.100.100.100", Name: "*", TTL: 600, Type: "A"},
					{Data: "100.100.100.100", Name: "@", TTL: 600, Type: "A"},
				}))),
				Header: make(http.Header),
			}
		} else if req.Method == http.MethodPut {
			assert.Equal(t, req.URL.String(), "localhost:3333/v1/domains/example.com/records/A")
			assert.Equal(t, req.Method, http.MethodPut)
			b, err := ioutil.ReadAll(req.Body)
			assert.NoError(t, err)
			assert.Equal(t, string(b), recordStrings([]Record{
				{Data: "101.101.101.101", Name: "*", TTL: 600, Type: "A"},
				{Data: "101.101.101.101", Name: "@", TTL: 600, Type: "A"},
			}))
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
	report, err := runner.CheckAndUpdate("example.com", "101.101.101.101", WithEndpoint("localhost:3333"), WithHttpClient{Client: client})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, report.DidUpdate, true)
}

func TestItUpdatesAllFoundNames(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method == http.MethodGet {
			return &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(recordStrings([]Record{
					{Data: "100.100.100.100", Name: "*", TTL: 600, Type: "A"},
					{Data: "100.100.100.100", Name: "@", TTL: 600, Type: "A"},
					{Data: "100.100.100.100", Name: "bob", TTL: 600, Type: "A"},
				}))),
				Header: make(http.Header),
			}
		} else if req.Method == http.MethodPut {
			assert.Equal(t, req.URL.String(), "localhost:4444/v1/domains/example.com/records/A")
			assert.Equal(t, req.Method, http.MethodPut)
			b, err := ioutil.ReadAll(req.Body)
			assert.NoError(t, err)
			assert.Equal(t, string(b), recordStrings([]Record{
				{Data: "101.101.101.101", Name: "*", TTL: 600, Type: "A"},
				{Data: "101.101.101.101", Name: "@", TTL: 600, Type: "A"},
				{Data: "101.101.101.101", Name: "bob", TTL: 600, Type: "A"},
			}))

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
	_, err := runner.CheckAndUpdate("example.com", "101.101.101.101", WithEndpoint("localhost:4444"), WithHttpClient{Client: client})
	if err != nil {
		t.Fatal(err)
	}
}

func TestItUpdatesOnlySpecifiedRecordsWhenSpecified(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method == http.MethodGet {
			return &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(recordStrings([]Record{
					{Data: "100.100.100.100", Name: "*", TTL: 600, Type: "A"},
					{Data: "100.100.100.100", Name: "@", TTL: 600, Type: "A"},
					{Data: "100.100.100.100", Name: "bob", TTL: 600, Type: "A"},
				}))),
				Header: make(http.Header),
			}
		} else if req.Method == http.MethodPut {
			assert.Equal(t, req.URL.String(), "localhost:3333/v1/domains/asdf.com/records/A")
			assert.Equal(t, req.Method, http.MethodPut)
			b, err := ioutil.ReadAll(req.Body)
			assert.NoError(t, err)
			assert.Equal(t, string(b), recordStrings([]Record{
				{Data: "101.101.101.101", Name: "*", TTL: 600, Type: "A"},
				{Data: "100.100.100.100", Name: "@", TTL: 600, Type: "A"},
				{Data: "100.100.100.100", Name: "bob", TTL: 600, Type: "A"},
			}))

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
	_, err := runner.CheckAndUpdate("asdf.com", "101.101.101.101", WithRecordName("*"), WithEndpoint("localhost:3333"), WithHttpClient{Client: client})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDryRunReturnsWhatWillBeUpdated(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method == http.MethodGet {
			return &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(recordStrings([]Record{
					{Data: "100.100.100.100", Name: "*", TTL: 600, Type: "A"},
					{Data: "100.100.100.100", Name: "@", TTL: 600, Type: "A"},
				}))),
				Header: make(http.Header),
			}
		} else if req.Method == http.MethodPut {
			assert.Fail(t, "put should not have been called")
			return nil
		} else {
			return &http.Response{
				StatusCode: 500,
				Header:     make(http.Header),
			}
		}
	})

	runner := &Updater{}
	report, err := runner.CheckAndUpdate("asdf.com", "101.101.101.101", WithRecordName("*"), WithEndpoint("localhost:3333"), WithHttpClient{Client: client}, WithDryRun(true))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(report.Records), 2)
	assert.Equal(t, report.DidUpdate, false)
}

func TestItReportsFAilureWhenItFailsToUpdate(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 500,
			Header:     make(http.Header),
		}
	})

	runner := &Updater{}
	_, err := runner.CheckAndUpdate("asdf.com", "101.101.101.101", WithEndpoint("localhost:3333"), WithHttpClient{Client: client})
	assert.Error(t, err)
}

func TestItReportsNoChangeWhenNoChangeIsNeeded(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method == http.MethodGet {
			return &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(recordStrings([]Record{
					{Data: "100.100.100.100", Name: "*", TTL: 600, Type: "A"},
					{Data: "100.100.100.100", Name: "@", TTL: 600, Type: "A"},
				}))),
				Header: make(http.Header),
			}
		} else if req.Method == http.MethodPut {
			assert.Fail(t, "put should not have been called")
			return nil
		} else {
			return &http.Response{
				StatusCode: 500,
				Header:     make(http.Header),
			}
		}
	})

	runner := &Updater{}
	report, err := runner.CheckAndUpdate("asdf.com", "100.100.100.100", WithRecordName("*"), WithEndpoint("localhost:3333"), WithHttpClient{Client: client})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, report.DidUpdate, false)
}
