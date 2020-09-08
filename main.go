package main

import (
    "net/http"
	"fmt"
	"log"
)

func main() {
// baseurl:= "https://www.youtube.com/watch?v=u8_JTzSSOIM"

http.HandleFunc("/", getUserUrl)
if err := http.ListenAndServe(":8080", nil); err != nil {
	log.Fatal(err)
}
}

// Handles the user input.
func getUserUrl(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
	}
	switch r.Method {
		case "GET":     
			 http.ServeFile(w, r, "home.html")
		case "POST":
			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			err := r.ParseForm(); 
			checkError(err, w)
			
			url := r.FormValue("url")
			getUrlInfo(url, w)
			// fmt.Fprintf(w, "Url = %s\n", url)
		default:
			fmt.Fprintf(w, "GET and POST methods are supported.")
		}
}
// Handles retrieving information of page.
func getUrlInfo(url string, w http.ResponseWriter) {
	response, err := http.Get(url)
	if !checkError(err, w) {
	fmt.Fprintf(w, "Response = %s\n", response)
	}
}
// Handles assigning the values from the info that getUrlInfo retrieves to result.html .
func createResponsePage(){

}
// Checks for errors from Http request
func checkError(err error, w http.ResponseWriter) bool{
	if err != nil {
		fmt.Fprintf(w, "Error = %s\n", err)
		return true;
	}
	return false;
}