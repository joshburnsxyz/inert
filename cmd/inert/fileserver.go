package main

import (
	"errors"
	"io"
	"net/http"
)

func makeFS(dir string) (http.HandlerFunc, error) {
	var (
		err      error
		handler  http.HandlerFunc
		htmlHead string = "<head>\n<link rel=\"stylesheet\" type=\"text/css\" href=\"https://cdnjs.cloudflare.com/ajax/libs/bulma/1.0.0/css/bulma.min.css\"/></head>"
	)

	fs := http.FileServer(http.Dir(dir))

	handler = func(w http.ResponseWriter, r *http.Request) {
		var (
			url   = r.URL.Path
			isDir = url[len(url)-1] == '/'
		)

		if isDir {
			err = nil
			io.WriteString(w, "<!doctype HTML>")
			io.WriteString(w, "\n")
			io.WriteString(w, "<html data-theme=\"dark\">")
			io.WriteString(w, "\n")
			io.WriteString(w, htmlHead)
			io.WriteString(w, "\n")
			io.WriteString(w, "<body>")
			io.WriteString(w, "\n")
			fs.ServeHTTP(w, r)
			io.WriteString(w, "\n")
			io.WriteString(w, "</body>")
			io.WriteString(w, "\n")
			io.WriteString(w, "</html>")
		} else {
			err = errors.New("Directory not found, could not generate HTTP handler.")
		}

	}

	return handler, err
}
