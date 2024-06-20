package main

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type fileRecord struct {
	Name string
	Size int64
	Path string
}

type customError struct {
	message string
	code    int
}

func (e *customError) Error() string {
	return e.message
}

func buildFileRecord(d os.DirEntry) (fileRecord, error) {
	info, err := d.Info()
	if err != nil {
		return fileRecord{}, &customError{message: "failed to get file info", code: http.StatusInternalServerError}
	}
	return fileRecord{
		Name: d.Name(),
		Size: info.Size(),
		Path: d.Name(),
	}, nil
}

func calculateAbsolutePath(dir, path string) string {
	return filepath.Join(dir, path)
}

func makeFS(dir string) (http.HandlerFunc, error) {

	// Build HTML template
	tmpl, err := template.New("file_index").Parse(`
<!DOCTYPE html>
<html data-theme="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.PageTitle}}</title>
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
        <h1 class="title">{{.PageTitle}}</h1>
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
		return nil, &customError{message: "failed to parse template", code: http.StatusInternalServerError}
	}

	// Build HTTP handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		filePath := calculateAbsolutePath(dir, r.URL.Path)
		file, err := os.Stat(filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pageTitle, _ := filepath.Abs(dir)
		data := struct {
			IsDir      bool
			Error      error
			DirEntries []fileRecord
			PageTitle  string
		}{
			IsDir:      file.IsDir(),
			Error:      nil,
			DirEntries: nil,
			PageTitle:  pageTitle,
		}

		if file.IsDir() {
			dirEntries, err := os.ReadDir(filePath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var files []fileRecord
			for _, d := range dirEntries {
				record, err := buildFileRecord(d)
				if err != nil {
					http.Error(w, err.Error(), err.(*customError).code)
					return
				}
				files = append(files, record)
			}

			data.DirEntries = files
			err = tmpl.ExecuteTemplate(w, "file_index.html", data)
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
