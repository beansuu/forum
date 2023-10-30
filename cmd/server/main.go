package main

import (
	"fmt"
	"forum/internal/database"
	"forum/internal/posts"
	"html/template"
	"net/http"
	"path"
	"strconv"

	"forum/internal/users"

	"golang.org/x/crypto/bcrypt"
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
	c, err := r.Cookie("session_id")
	if err != nil || c == nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	userIDStr, valid := users.GetUserIDFromSession(c.Value)
	if !valid {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	// Convert userIDStr (which is a string) to int64
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error converting userIDStr to int64:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Fetch all posts for this user
	userPosts, err := posts.GetAllPostsForUser(userID)
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}

	templates.ExecuteTemplate(w, "all-posts.html", userPosts)
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

func handleRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

		// Check if email or username already exists
		existingUser, err := users.GetUserByEmailOrUsername(username)
		if err != nil {
			http.Error(w, "Error checking user", http.StatusInternalServerError)
			return
		}
		if existingUser != nil {
			http.Error(w, "Username or Email already exists", http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		_, err = users.CreateUser(username, string(hashedPassword), email)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		// Render the registration page template
		templates.ExecuteTemplate(w, "register.html", nil)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		emailOrUsername := r.FormValue("emailOrUsername")
		password := r.FormValue("password")

		// Check if user exists
		user, err := users.GetUserByEmailOrUsername(emailOrUsername)
		if err != nil {
			http.Error(w, "User doesn't exist", http.StatusUnauthorized)
			return
		}

		// Verify the password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Error(w, "Wrong password", http.StatusUnauthorized)
			return
		}

		// Here, handle setting session and cookie logic to log the user in

		http.Redirect(w, r, "/all-posts", http.StatusSeeOther)
	} else {
		// Render the login page template
		templates.ExecuteTemplate(w, "login.html", nil)
	}
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
	http.HandleFunc("/register", handleRegistration)
	http.HandleFunc("/login", handleLogin)

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
