package main

import (
    "net/http"
	"fmt"
	"log"
	"errors"
	"net/url"
	"html/template"
	"io/ioutil"
	"strings"

)

//URL is a global since we want to access it throughout in different handler methods
var userInput string

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

			// If there are no errors we set the url to the user inputß
			userInput = r.FormValue("url")

			// Validating user input to be a real url
			u, err := url.ParseRequestURI(userInput)
			_ = u
			if err != nil {
				var er = errors.New("Please enter a valid URL.")
				checkError(er, w, r)
				return	
			}

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
// Analyzes webpage and creates the data object that is displayed in result page.

func analyzePage(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get(userInput)
	if err != nil {
		var er = errors.New("Cannot get info from URL.")
		checkError(er, w, r)
		return	
	}
	// Getting html body
	body, err:= ioutil.ReadAll(response.Body)
    pageContent := string(body)

	// Getting page title
	pageTitle:= getPageTitle(pageContent)
	// Checking if login exists
	loginExists:=getLoginExists(pageContent)

	// Print out the result
	fmt.Printf("Page title: %s\n", pageTitle)

	data := map[string]interface{}{"htmlVersion": "012031023021", "pageTitle": pageTitle, "loginForm" : loginExists}
	outputHTML(w, "static/result.html", data)
	
}
func getLoginExists(pageContent string) bool{

}
func getPageTitle(pageContent string) string{
	// Find a substr
	titleStartIndex := strings.Index(pageContent, "<title>")
	if titleStartIndex == -1 {
		fmt.Println("No title element found")
		return "N/A"
	}
	// The start index of the title is the index of the first
	// character, the < symbol. We don't want to include
	// <title> as part of the final value, so let's offset
	// the index by the number of characers in <title>
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
func checkError(err error, w http.ResponseWriter, r *http.Request) bool{
	if err != nil {
		errorMessage := map[string]interface{}{"errorMessage": err}
    	outputHTML(w, "static/errorPage.html", errorMessage)
		return true;
	}
	return false;
}