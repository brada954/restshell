package shell

import (
	"net/http"
	"net/url"
	"testing"
)

func TestNewQueryParamAuthNoParams(t *testing.T) {
	var startingUrl = "http://abc.com/?x=123"

	auth := NewQueryParamAuth()
	if len(auth.KeyPairs) != 0 {
		t.Errorf("Expected zero parameters; found %d", len(auth.KeyPairs))
	}

	req := &http.Request{}
	var err error
	req.URL, err = url.Parse(startingUrl)
	if err != nil {
		t.Errorf("Unexpected error parsing starting url: %s", err.Error())
	}

	auth.AddAuth(req)
	if req.URL.String() != startingUrl {
		t.Errorf("Unexpected url; %s!=%s", startingUrl, req.URL.String())
	}
}

func TestNewQueryParamAuthWithParams(t *testing.T) {
	var startingUrl = "http://abc.com/?x=123"
	var expectedUrl = "http://abc.com/?x=123&MORID=MYID&MORKEY=MYSECRET&HASH=Test+Space"

	auth := NewQueryParamAuth("MORID", "MYID", "MORKEY", "MYSECRET", "HASH", "Test Space")
	if len(auth.KeyPairs) != 3 {
		t.Errorf("Expected zero parameters; found %d", len(auth.KeyPairs))
	}

	req := &http.Request{}
	var err error
	req.URL, err = url.Parse(startingUrl)
	if err != nil {
		t.Errorf("Unexpected error parsing starting url: %s", err.Error())
	}

	auth.AddAuth(req)
	if req.URL.String() != expectedUrl {
		t.Errorf("Unexpected url; %s!=%s", expectedUrl, req.URL.String())
	}
}
