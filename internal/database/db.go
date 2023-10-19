package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Initialize() error {
	var err error
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}

	createPostsTable := `
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			author TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err = db.Exec(createPostsTable)
	if err != nil {
		return err
	}

	createCommentsTable := `
	CREATE TABLE comments (
		id INTEGER PRIMARY KEY,
		post_id INTEGER,
		author TEXT,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(post_id) REFERENCES posts(id)
	);
	`
	_, err = db.Exec(createCommentsTable)
	if err != nil {
		return err
	}

	createLikesTable := `
	CREATE TABLE post_likes (
		id INTEGER PRIMARY KEY,
		post_id INTEGER,
		user_id INTEGER,
		like_value INTEGER,  -- +1 for like, -1 for dislike
		FOREIGN KEY(post_id) REFERENCES posts(id)
	);
	
	`
	_, err = db.Exec(createLikesTable)
	if err != nil {
		return err
	}

	createUsersTable := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY,
		username TEXT UNIQUE,
		password TEXT,  -- preferably a hashed version of the password
		email TEXT UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	`
	_, err = db.Exec(createUsersTable)
	if err != nil {
		return err
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}
