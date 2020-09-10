package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/steelx/extractlinks"
)

//URL is a global since we want to access it throughout in different handler methods
var userInput string
var mainURL string
var hostURL string
var domainName string

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/result", resultHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

// Handles the user input.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		var ErrNotFound = errors.New("Only the home page is accessible.")
		checkError(ErrNotFound, w, r)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "static/home.html")
	case "POST":

		// If there are no errors we set the url to the user input
		userInput = r.FormValue("url")

		// Validating user input to be a real url
		u, err := url.ParseRequestURI(userInput)
		_ = u
		if err != nil {
			var er = errors.New("Please enter a valid URL.")
			checkError(er, w, r)
			return
		}

		hostURL = u.Host
		mainURL = extractMainUrl(hostURL)
		domainName = extractDomainName(hostURL)

		// Redirecting to result page
		http.Redirect(w, r, "/result", 302)

	default:
		fmt.Fprintf(w, "GET and POST methods are supported.")
	}
}

// Handles retrieving results page.
func resultHandler(w http.ResponseWriter, r *http.Request) {

	// Checking if redirect was from home page, otherwise throw error message
	if r.Header.Get("Referer") == "http://localhost:8080/" {

		if r.Method == http.MethodGet {
			analyzePage(w, r)
		}
	} else {
		var er = errors.New("Cannot access page directly.")
		checkError(er, w, r)
		return
	}
}

// Returns an array with index: level and value: count of headings at level
func getHeadingsCount() {

}

// Analyzes webpage and creates the data object that is displayed in result page.

func analyzePage(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get(userInput)
	defer response.Body.Close()
	if err != nil {
		var er = errors.New("Cannot get info from URL.")
		checkError(er, w, r)
		return
	}

	// Getting html body
	body, err := ioutil.ReadAll(response.Body)
	pageContent := string(body)

	// Resetting response body
	response.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Extracting links using external package
	links, err := extractlinks.All(response.Body)
	// Calculating link info
	internalLinksCount := getInternalLinksCount(links)
	externalLinksCount := len(links) - internalLinksCount
	inaccessibleLinksCount := getInaccessibleLinks(links)
	// Getting html version
	htmlVersion := getHTMLVersion(pageContent)
	// Getting page title
	pageTitle := getPageTitle(pageContent)
	// Checking if login exists
	loginExists := getLoginExists(pageContent)

	data := map[string]interface{}{"htmlVersion": htmlVersion, "pageTitle": pageTitle, "loginForm": loginExists, "internalLinksCount": internalLinksCount, "externalLinksCount": externalLinksCount, "inaccessibleLinksCount": inaccessibleLinksCount, "url": userInput}
	outputHTML(w, "static/result.html", data)

}
func extractMainUrl(content string) string {

	startIndex := 0

	endIndex := strings.Index(content, ".")
	if content[startIndex:endIndex] == "www" {
		startIndex += 4
		slicedWWW := (content[startIndex:len(content)])
		endIndex = strings.Index(slicedWWW, ".")

		slicedWWW = (slicedWWW[0:endIndex])

		return slicedWWW
	}

	mainURLRet := (content[startIndex:endIndex])

	return mainURLRet
}
func extractDomainName(content string) string {
	startIndex := 0

	endIndex := strings.Index(content, ".")

	if content[startIndex:endIndex] == "www" {
		startIndex += 4
		slicedWWW := (content[startIndex:len(content)])
		endIndex = strings.Index(slicedWWW, ".")
		slicedWWW = (content[endIndex+1 : len(content)])
		return slicedWWW
	}

	mainURLRet := (content[endIndex+1 : len(content)])
	return mainURLRet
}

//	Returns inaccessible links of current website
func getInaccessibleLinks(links []extractlinks.Link) int {
	var count int
	for i, link := range links {
		_ = i
		currentLink := link.Href

		if strings.Index(currentLink, ".") != 0 {

			// In case of a url that is not defined but is on the home url
			if strings.Index(currentLink, "/") == 0 && !strings.Contains(currentLink, "http") {
				currentLink = hostURL + currentLink
			}
			// http.Get only works with http or https before hostname
			if !strings.Contains(currentLink, "http") {
				currentLink = "http://" + currentLink
			}
			resp, err := http.Get(currentLink)
			if err != nil {
				count++
				log.Fatal(err)
				continue
			}
			if resp.StatusCode != 200 {
				count++
			}
		}
	}
	fmt.Println(count)
	return count
}

//	Returns internal links of current website
func getInternalLinksCount(links []extractlinks.Link) int {
	var count int
	for i, link := range links {
		_ = i
		currentLink := link.Href
		if strings.Contains(currentLink, mainURL) || strings.Index(currentLink, "/") == 0 {
			count++
		}
	}
	return count
}

func getHTMLVersion(pageContent string) string {

	// Find the beginning html tag of input and store
	docStartIndex := strings.Index(pageContent, "<!DOCTYPE")
	if docStartIndex == -1 {
		fmt.Println("No doc tag found")
		return "Latest"
	}
	// We don't want to include
	// <DOCTYPE> as part of the final value, so we offset
	// the index by the number of characters in <DOCTYPE>
	docStartIndex += 9

	// Find the index of the closing tag
	docEndIndex := strings.Index(pageContent, ">")
	if docEndIndex == -1 {
		fmt.Println("No doc end tag  found")
		return "Latest"
	}

	docContent := (pageContent[docStartIndex:docEndIndex])
	fmt.Println(docContent)
	versionLookup := strings.Index(docContent, "HTML")
	if versionLookup == -1 {
		fmt.Println("No HTML version found")
		return "Latest"
	}
	version := (pageContent[versionLookup : versionLookup+5])
	return version
}
func getLoginExists(pageContent string) string {
	// Find the beginning html tag of input and store
	formStartIndex := strings.Index(pageContent, "<form")
	if formStartIndex == -1 {
		fmt.Println("No form tag found")
		return "No"
	}
	// We don't want to include
	// <form> as part of the final value, so we offset
	// the index by the number of characters in <form>
	formStartIndex += 5

	// Find the index of the closing tag
	formEndIndex := strings.Index(pageContent, "</form>")
	if formEndIndex == -1 {
		fmt.Println("No form tag end found")
		return "No"
	}

	formContent := (pageContent[formStartIndex:formEndIndex])

	pwdLookup := strings.Index(formContent, "password")
	if pwdLookup == -1 {
		fmt.Println("No login form found")
		return "No"
	}
	return "Yes"
}

func getPageTitle(pageContent string) string {
	// Find the beginning html tag of title and store
	titleStartIndex := strings.Index(pageContent, "<title>")
	if titleStartIndex == -1 {
		fmt.Println("No title element found")
		return "N/A"
	}
	// We don't want to include
	// <title> as part of the final value, so we offset
	// the index by the number of characters in <title>
	titleStartIndex += 7

	// Find the index of the closing tag
	titleEndIndex := strings.Index(pageContent, "</title>")
	if titleEndIndex == -1 {
		fmt.Println("No closing tag for title found for Title.")
		return "N/A"
	}

	// Copy the substring in to a separate variable so the
	// variables with the full document data can be garbage collected
	pageTitle := []byte(pageContent[titleStartIndex:titleEndIndex])
	return string(pageTitle)
}

// Handles assigning the values from the info that resultHandler retrieves to result.html .
func outputHTML(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Checks for errors from Http request
func checkError(err error, w http.ResponseWriter, r *http.Request) bool {
	if err != nil {
		errorMessage := map[string]interface{}{"errorMessage": err}
		outputHTML(w, "static/errorPage.html", errorMessage)
		return true
	}
	return false
}
