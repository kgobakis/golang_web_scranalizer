package main

import (
	"testing"
)

func TestExtractMainUrl(t *testing.T) {
	url1 := "www.hotmail.com"
	url2 := "hotmail.net"

	if extractMainUrl(url1) != "hotmail" || extractMainUrl(url2) != "hotmail" {
		t.Error("Expected hotmail.")
	}

}
func TestExtractDomainName(t *testing.T) {
	url1 := "hotmail.com"
	url2 := "gmail.com"

	if extractDomainName(url1) != "com" || extractDomainName(url2) != "com" {
		t.Error("Expected hotmail.")
	}

}
