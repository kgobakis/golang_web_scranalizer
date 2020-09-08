package main

import (
    "net/http"
	"fmt"
	"os"
	// "log"
)

func main() {
// baseurl:= "https://www.youtube.com/watch?v=u8_JTzSSOIM"

http.HandleFunc("/", getUserUrl)
err := http.ListenAndServe(":8080", nil);
checkError(err)
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
			checkError(err)
			
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
	checkError(err)
	fmt.Fprintf(w, "Response = %s\n", response)

}
// Checks for errors from Http request
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
    }
}