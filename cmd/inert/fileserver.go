package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

// fileRecord represents a file or directory entry.
type fileRecord struct {
	Name  string
	Size  int64
	Path  string
	IsDir bool
}

// customError represents a custom error type.
type customError struct {
	message string
	code    int
}

// Error returns the error message.
func (e *customError) Error() string {
	return e.message
}

// buildFileRecord builds a fileRecord from an os.DirEntry.
// Returns a customError if there is an error getting file info.
func buildFileRecord(dirPath string, d os.DirEntry) (fileRecord, error) {
	info, err := d.Info()
	if err != nil {
		return fileRecord{}, &customError{message: "failed to get file info", code: http.StatusInternalServerError}
	}
	return fileRecord{
		Name:  d.Name(),
		Size:  info.Size(),
		Path:  filepath.Join(dirPath, d.Name()),
		IsDir: info.IsDir(),
	}, nil
}

// calculateAbsolutePath calculates the absolute path of a file or directory.
func calculateAbsolutePath(dir, path string) string {
	return filepath.Join(dir, path)
}

// makeFS creates an HTTP handler that serves a directory listing.
// The handler builds an HTML template and returns it as a response.
// If the requested path is a directory, it lists its contents.
// If the requested path is a file, it serves the file.
// Returns a customError if there is an error parsing the template or getting file info.
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
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.3/css/all.min.css">
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
	.file-icon .folder-icon {
	    display: inline-block;
	    vertical-align: middle;
	    margin-right: 5px;
	}
        .file-icon {
            color: #2196F3;
        }
        .folder-icon {
            color: #FFC107;
        }
    </style>
</head>
<body>
<div class="container">
    {{if .IsDir}}
        <h1 class="title">{{.PageTitle}}</h1>
        <ul class="menu-list">
            {{range .DirEntries}}
		<a href={{.Path}}>
		    <li>
			<div class="columns">
			    <div class="column">
				{{if .IsDir}}
				    <i class="fas fa-folder folder-icon"></i>
				{{else}}
				    <i class="fas fa-file file-icon"></i>
				{{end}}
			    </div>
			    <div class="column is-four-fifths">
				{{.Name}}
			    </div>
			    <div class="column">
				{{if .IsDir}}
				{{else}}
				    <span>{{.Size}} bytes</span>
				{{end}}
			    </div>
			</div>
		    </li>
		</a>
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
				fmt.Println(r.URL)
				record, err := buildFileRecord(r.URL.Path, d)
				if err != nil {
					http.Error(w, err.Error(), err.(*customError).code)
					return
				}
				files = append(files, record)
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
