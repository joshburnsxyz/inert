package main

import (
	"io"
	"net/http"
	"errors"
)

func makeFS(dir string) (http.HandlerFunc, error) {
	var (
	    htmlHead string = "<head>\n<link rel=\"stylesheet\" type=\"text/css\" href=\"https://raw.githubusercontent.com/raj457036/attriCSS/master/themes/midnight-green.css\"/>\n</head>"
	    err error
	    handler http.HandlerFunc
	)

	fs := http.FileServer(http.Dir(dir))

	handler = func(w http.ResponseWriter, r *http.Request) {
		var (
			url = r.URL.Path
			isDir = url[len(url)-1] == '/'
		)

		// Create HTTP Server
		
		// If directory is found generate the page,
		// and return the handler
		if isDir {
			err = nil
			//w.Header().Set("Conent-Type", "text/html")
			io.WriteString(w, "<!doctype HTML>")
			io.WriteString(w, "\n")
			io.WriteString(w, "<html>")
			io.WriteString(w, htmlHead)
			io.WriteString(w, "\n")
			io.WriteString(w, "<body>")
			fs.ServeHTTP(w, r)
			io.WriteString(w, "</body>")
			io.WriteString(w, "</html>")
		} else {
			err = errors.New("Directory not found, could not generate HTTP handler.")
		}

	}
	return handler,err
}
