package main

import (
  "testing"
  "regexp"
  //"net/http"
  //"net/http/httptest"
)

/*
type TinyURL struct {
  ID string
  LongURL string
  ShortURL string
}
*/

var testTinyURL = TinyURL{"ZwvkBAv", "http://google.com", "http://localhost:8000/dZwvkBAv"}

func TestCheckHTTP(t *testing.T) {
  var urlWithHTTP = CheckHTTP("http://google.com")
  var urlWithHTTPS = CheckHTTP("https://google.com")
  var urlNoHTTP = CheckHTTP("google.com")
  var noUrl = CheckHTTP("")

  if urlWithHTTP != "http://google.com"{
    t.Error("Expected http://google.com, got ", urlWithHTTP)
  }
  if urlWithHTTPS != "https://google.com"{
    t.Error("Expected https://google.com, got ", urlWithHTTPS)
  }
  if urlNoHTTP != "http://google.com"{
    t.Error("Expected http://google.com, got ", urlNoHTTP)
  }
  if noUrl != ""{
    t.Error("Expected \"\", got ", noUrl)
  }
}

func TestCreateURL(t *testing.T) {
  var tempTinyURL = CreateURL(testTinyURL)
  matched, _ := regexp.MatchString("http://localhost:8000/d" + tempTinyURL.ID, tempTinyURL.ShortURL)
  if !(matched) {
    t.Error("Expected the ID", tempTinyURL.ID, "to match the end of the ShortURL", tempTinyURL.ShortURL)
  }
}
