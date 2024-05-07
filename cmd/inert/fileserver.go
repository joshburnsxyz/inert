package main

import (
	"net/http"
	"html/template"
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
<!doctype HTML>
<html data-theme="dark">
<head>
    <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/1.0.0/css/bulma.min.css"/>
</head>
<body>
{{if .IsDir}}
<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Size</th>
        </tr>
    </thead>
    <tbody>
        {{range .DirEntries}}
        <tr>
            <td><a href="{{.Name}}">{{.Name}}</a></td>
            <td>{{.Size}}</td>
        </tr>
        {{end}}
    </tbody>
</table>
{{else}}
<h1>Error</h1>
<p>{{.Error}}</p>
{{end}}
</body>
</html>
`)
	if err != nil {
		return nil,err
	}

	// Build HTTP handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		file,err := os.Stat(dir + r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Build struct to hold data to populate the HTML template with.
		data := struct {
			IsDir bool
			Error error
			DirEntries []fileRecord
		}{
			IsDir: file.IsDir(),
			Error: nil,
			DirEntries: nil,
		}

		if file.IsDir() {
			data.Error = nil
			dirEntries, err := os.ReadDir(dir + r.URL.Path)
			var files []fileRecord

			// Build a custom data struct so we can calculate the absolute path to the file
			for _,d := range dirEntries {
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
			fs.ServeHTTP(w,r)
		}

	}

	// Upon success, return the handler and no error
	return handler,nil
}
