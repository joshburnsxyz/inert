package main

import (
	"io"
	"net/http"
	"errors"
)

func makeFS(dir string) (http.HandlerFunc, error) {
	fs := http.FileServer(http.Dir(dir))
	var handler http.HandlerFunc
	var cssLink string = "<link rel=\"stylesheet\" href=\"https://raw.githubusercontent.com/raj457036/attriCSS/master/themes/midnight-green.css\"/>"
	var err error
	handler = func(w http.ResponseWriter, r *http.Request) {
		var (
			url = r.URL.Path
			isDir = url[len(url)-1] == '/'
		)

		// Create HTTP Server
		fs.ServeHTTP(w, r)

		// If directory is found generate the page,
		// and return the handler
		if isDir {
			w.Header().Set("Conent-Type", "text/html")
			io.WriteString(w, cssLink)
			io.WriteString(w, "\n")
			err = nil
		} else {
			err = errors.New("Directory not found, could not generate HTTP handler.")
		}
	}
	return handler,err
}
