package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

type Filename struct {
	Value    string
	Selected bool
}

type FileSource struct {
	Path      string
	Filenames []Filename
}

type Data struct {
	Title       string
	Content     string
	FileSources []FileSource
}

func handleRequest(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/favicon.ico" {
		http.ServeFile(writer, request, "public/favicon.ico")
		return
	}

	data := Data{
		Title: config.Name,
	}

	dirs, err := os.ReadDir(storageDir)
	if err != nil {
		fmt.Fprint(writer, "Can't read storage dir")
		return
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			fileSource := FileSource{
				Path: dir.Name(),
			}

			sourceDir := storageDir + "/" + dir.Name()
			files, err := os.ReadDir(sourceDir)
			if err != nil {
				fmt.Fprint(writer, "Can't read dir "+sourceDir)
				return
			}
			for _, file := range files {
				if !file.IsDir() {
					filename := Filename{
						Value:    file.Name(),
						Selected: false,
					}
					if request.URL.Path == ("/" + fileSource.Path + "/" + filename.Value) {
						filename.Selected = true
					}
					fileSource.Filenames = append(fileSource.Filenames, filename)
				}
			}

			data.FileSources = append(data.FileSources, fileSource)
		}
	}

	bytes, err := os.ReadFile(storageDir + request.URL.Path)
	if err == nil {
		data.Content = string(bytes)
	} else {
		data.Content = "Unknown log file\n\n[Config]\n" +
			"App name: '" + config.Name + "'\n" +
			"Refresh interval: " + strconv.Itoa(config.Interval) + "s"
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(writer, data)
}
