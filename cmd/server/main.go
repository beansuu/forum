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

func displayAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := posts.GetAllPosts()
	if err != nil {
		fmt.Println("Error retrieving posts:", err) // Print to console
		http.Error(w, "Unable to retrieve posts", http.StatusInternalServerError)
		return
	}
	err = templates.ExecuteTemplate(w, "posts.html", posts)
	if err != nil {
		fmt.Println("Template execution error:", err) // Print to console
	}
}

func getPostByID(postID string) (*posts.Post, error) {
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		return nil, err
	}
	return posts.GetPost(id)
}

func handleLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := r.FormValue("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Fetch post from DB
	_, err = getPostByID(postIDStr)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Increment like count (and ensure the same user doesn't like more than once)
	err = incrementLike(postID)
	if err != nil {
		http.Error(w, "Error processing like", http.StatusInternalServerError)
		return
	}

	// Redirect back to the post or wherever you want
	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}

func handleDislike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := r.FormValue("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Fetch post from DB to ensure it exists
	_, err = getPostByID(postIDStr)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Increment dislike count (and ensure the same user doesn't dislike more than once)
	err = incrementDislike(postID)
	if err != nil {
		http.Error(w, "Error processing dislike", http.StatusInternalServerError)
		return
	}

	// Redirect back to the post or wherever you want
	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}

func incrementLike(postID int64) error {
	db := database.GetDB()
	_, err := db.Exec("UPDATE posts SET likes=likes+1 WHERE id=?", postID)
	if err != nil {
		return fmt.Errorf("failed to increment likes: %v", err)
	}
	return nil
}

func incrementDislike(postID int64) error {
	db := database.GetDB()
	_, err := db.Exec("UPDATE posts SET dislikes=dislikes+1 WHERE id=?", postID)
	if err != nil {
		return fmt.Errorf("failed to increment dislikes: %v", err)
	}
	return nil
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
	http.HandleFunc("/all-posts", displayAllPostsHandler)
	http.HandleFunc("/like", handleLike)
	http.HandleFunc("/dislike", handleDislike)

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
