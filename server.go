package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

//go:embed tmpl/*
var templateFS embed.FS

func main() {
	r := mux.NewRouter()

	// Static files (CSS, JS, images)
	staticFS, err := fs.Sub(templateFS, "tmpl/static")
	if err != nil {
		log.Fatal(err)
	}

	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/",
			http.FileServer(http.FS(staticFS)),
		))

	// Home page
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := templateFS.ReadFile("tmpl/index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	// Dynamic pages: /pages/about, /pages/menu, etc.
	r.HandleFunc("/pages/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		file := path.Join("tmpl/static/pages", name+".html")

		data, err := templateFS.ReadFile(file)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	log.Println("http://localhost:8000/")
	log.Fatal(http.ListenAndServe(":8000", r))
}
