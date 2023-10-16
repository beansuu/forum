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

func UpdatePost(id int64, title, content, author string) error {
	_, err := db.Exec("UPDATE posts SET title=?, content=?, author=? WHERE id=?", title, content, author, id)
	return err
}

func DeletePost(id int64) error {
	_, err := db.Exec("DELETE FROM posts WHERE id=?", id)
	return err
}
