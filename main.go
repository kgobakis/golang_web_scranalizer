package main

import (
    "net/http"
	"fmt"
	"log"
	"errors"
	"net/url"
"html/template"

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
			response, err := http.Get(userInput)
			_ = response
			if !checkError(err, w, r) {
				http.ServeFile(w, r, "static/result.html")
				// fmt.Fprintf(w, "Response = %s\n", response)
			}
		} 
   } else {
	var er = errors.New("Cannot access page directly.")
	checkError(er, w, r)
	return	
   }
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