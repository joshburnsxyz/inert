package main

import (
	"html/template"
	"net/http"
	"os"
)

type fileRecord struct {
	Name string
	Size int64
	Path string
}

func makeFS(dir string) (http.HandlerFunc, error) {

	// Build HTML template
	tmpl, err := template.New("file_index").Parse(`
<!DOCTYPE html>
<html data-theme="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Directory Listing</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/1.0.0/css/bulma.min.css">
    <style>
        body {
            padding-top: 40px;
            padding-bottom: 40px;
        }
        .container {
            max-width: 600px;
            margin: auto;
        }
        .error {
            color: red;
        }
    </style>
</head>
<body>
<div class="container">
    {{if .IsDir}}
        <h1 class="title">Directory Listing</h1>
        <ul class="menu-list">
            {{range .DirEntries}}
                <li><a href="{{.Name}}" class="has-text-white">{{.Name}}</a></li>
		<hr/>
            {{end}}
        </ul>
    {{else}}
        <h1 class="title">Error</h1>
        <p class="error">{{.Error}}</p>
    {{end}}
</div>
</body>
</html>
`)
	if err != nil {
		return nil, err
	}

	// Build HTTP handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Stat(dir + r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Build struct to hold data to populate the HTML template with.
		data := struct {
			IsDir      bool
			Error      error
			DirEntries []fileRecord
		}{
			IsDir:      file.IsDir(),
			Error:      nil,
			DirEntries: nil,
		}

		if file.IsDir() {
			data.Error = nil
			dirEntries, err := os.ReadDir(dir + r.URL.Path)
			var files []fileRecord

			// Build a custom data struct so we can calculate the absolute path to the file
			for _, d := range dirEntries {
				info, _ := d.Info()
				new_file := fileRecord{
					Name: d.Name(),
					Size: info.Size(),
					Path: d.Name(),
				}
				files = append(files, new_file)

			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data.DirEntries = files
			err = tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			fs := http.FileServer(http.Dir(dir))
			fs.ServeHTTP(w, r)
		}

	}

	// Upon success, return the handler and no error
	return handler, nil
}
