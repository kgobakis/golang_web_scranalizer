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

func TestGetLoginExists(t *testing.T) {
	test1 := "<html/> <h1/> <form password/>"
	test2 := "<html/> <h1/> <form password> </form>"

	if getLoginExists(test1) != "No" {
		t.Error("Expected No.")
	}

	if getLoginExists(test2) != "Yes" {
		t.Error("Expected Yes.")
	}
}
func TestGetPageTitle(t *testing.T) {
	test1 := "<title/> Hi <title>"
	test2 := "<title> Hello</title>"

	if getPageTitle(test1) != "N/A" {
		t.Error("Expected No.")
	}

	if getPageTitle(test2) != " Hello" {
		t.Error("Expected Yes.")
	}
}
