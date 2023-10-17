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

func updateFormHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("postID")
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
	templates.ExecuteTemplate(w, "update.html", post)
}

func main() {

	err := database.Initialize()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the forum!")
	})

	http.HandleFunc("/update-form", updateFormHandler)
	http.HandleFunc("/create-post", createPostHandler)
	http.HandleFunc("/display-post", displayPostHandler)

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		// For demonstration, fetching values from query parameters.
		// Ideally, you'd use a form's POST data.
		postID, _ := strconv.ParseInt(r.URL.Query().Get("postID"), 10, 64)
		title := r.URL.Query().Get("title")
		content := r.URL.Query().Get("content")

		err := posts.UpdatePost(postID, title, content)
		if err != nil {
			http.Error(w, "Unable to update post", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Post updated successfully!")
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		postID, _ := strconv.ParseInt(r.URL.Query().Get("postID"), 10, 64)

		err := posts.DeletePost(postID)
		if err != nil {
			http.Error(w, "Unable to delete post", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Post deleted successfully!")
	})

	http.ListenAndServe(":8080", nil)
}
