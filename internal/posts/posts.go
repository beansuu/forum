package posts

import (
	"database/sql"
	"forum/internal/database"
)

type Post struct {
	ID      int64
	Title   string
	Content string
	Author  string
}

var db *sql.DB

func init() {
	db = database.GetDB()
}

func CreatePost(title, content, author string) (int64, error) {
	statement, err := db.Prepare("INSERT INTO posts(title, content, author) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(title, content, author)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func GetPost(id int64) (*Post, error) {
	post := &Post{}
	err := db.QueryRow("SELECT id, title, content, author FROM posts WHERE id=?", id).Scan(&post.ID, &post.Title, &post.Content, &post.Author)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func UpdatePost(postID int64, title string, content string) error {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}
	defer db.Close()

	statement, err := db.Prepare("UPDATE posts SET title=?, content=? WHERE id=?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(title, content, postID)
	if err != nil {
		return err
	}
	return nil
}

func DeletePost(postID int64) error {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}
	defer db.Close()

	statement, err := db.Prepare("DELETE FROM posts WHERE id=?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(postID)
	if err != nil {
		return err
	}
	return nil
}
