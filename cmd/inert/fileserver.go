package main

import (
	"io"
	"net/http"
)

func makeFS(dir string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(dir))
	var handler http.HandlerFunc
	var cssLink string = "<link rel=\"stylesheet\" href=\"https://raw.githubusercontent.com/raj457036/attriCSS/master/themes/midnight-green.css\"/>"
	
	handler = func(w http.ResponseWriter, r *http.Request) {
		var (
			url = r.URL.Path
			isDir = url[len(url)-1] == '/'
		)
		fs.ServeHTTP(w, r)
		if isDir {
			io.WriteString(w, cssLink)
		}
	}
	return handler
}
