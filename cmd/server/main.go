package main

import (
	"fmt"
	"forum/internal/database"
	"forum/internal/posts"
	"html/template"
	"net/http"
	"path"
	"strconv"
)

var templates *template.Template

func init() {
	// Parsing all templates from the web/templates directory
	templates = template.Must(template.ParseGlob(path.Join("web", "templates", "*.html")))
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// handling post creation
		title := r.FormValue("title")
		content := r.FormValue("content")
		author := r.FormValue("author")
		_, err := posts.CreatePost(title, content, author)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	templates.ExecuteTemplate(w, "create.html", nil)
	return
}

func displayPostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	post, err := posts.GetPost(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templates.ExecuteTemplate(w, "display.html", post)
}

func main() {

	err := database.Initialize()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the forum!")
	})

	http.HandleFunc("/create-post", createPostHandler)
	http.HandleFunc("/display-post", displayPostHandler)

	http.ListenAndServe(":8080", nil)
}
