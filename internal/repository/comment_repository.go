package repository

import (
	"database/sql"
	"fmt"
	"forum/internal/models"
)

// Comment is an interface that defines methods for interacting with comments in the database.
type Comment interface {
	CreateComment(comment *models.Comment) error
	GetComments(postID int) ([]*models.Comment, error)
	GetCommentByID(commentID int) (models.Comment, error)
	CommentHasLike(commentID int, username string) error
	CommentHasDislike(commentID int, username string) error
	RemoveLikeComment(commentID int, username string) error
	RemoveDislikeComment(commentID int, username string) error
	LikeComment(commentID int, username string) error
	DislikeComment(commentID int, username string) error
}

// CommentStorage is a struct that implements the Comment interface.
type CommentStorage struct {
	db *sql.DB
}

// NewCommentSqlite returns a new instance of CommentStorage.
func NewCommentSqlite(db *sql.DB) *CommentStorage {
	return &CommentStorage{db: db}
}

// CreateComment creates a new comment in the database.
func (c *CommentStorage) CreateComment(comment *models.Comment) error {
	query := fmt.Sprintf(`INSERT INTO comment (author, text, postid) values ($1, $2, $3)`)
	res, err := c.db.Exec(query, comment.Author, comment.Text, comment.PostID)
	if err != nil {
		return err
	}
	_, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("repository: create commentary: Insert query - %w", err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: create commentary: Insert query - %w", err)
	}

	return nil
}

// GetComments returns all comments for a given post ID.
func (c *CommentStorage) GetComments(postID int) ([]*models.Comment, error) {
	var comments []*models.Comment
	query := fmt.Sprintf(`SELECT * FROM comment WHERE postid = $1;`)
	rows, err := c.db.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("repository: get commentaries of the post: query - %w", err)
	}

	for rows.Next() {
		c := &models.Comment{}
		if err = rows.Scan(&c.ID, &c.Author, &c.PostID, &c.Text, &c.Likes, &c.DisLikes); err != nil {
			return nil, fmt.Errorf("repository: get commentaries of the post: query - %w", err)
		}
		comments = append(comments, c)
	}
	rows.Close()
	return comments, nil
}

// GetCommentByID returns a comment with a given ID.
func (c *CommentStorage) GetCommentByID(commentID int) (models.Comment, error) {
	var comment models.Comment

	query := `SELECT id, postid, author, text, like, dislike FROM comment WHERE id=$1;`
	row := c.db.QueryRow(query, commentID)

	err := row.Scan(&comment.ID, &comment.PostID, &comment.Author, &comment.Text, &comment.Likes, &comment.DisLikes)
	if err != nil {
		return models.Comment{}, fmt.Errorf("storage: get user by login: %w", err)
	}

	return comment, nil
}

// RemoveLikeComment removes a like from a comment.
func (s *CommentStorage) RemoveLikeComment(commentID int, username string) error {
	query := `DELETE FROM like WHERE commentId = $1 AND username = $2;`
	_, err := s.db.Exec(query, commentID, username)
	if err != nil {
		return fmt.Errorf("storage: remove like from comment: %w", err)
	}
	query = `UPDATE comment SET like = like - 1 WHERE id = $1;`
	_, err = s.db.Exec(query, commentID)
	if err != nil {
		return fmt.Errorf("storage: remove like from comment: %w", err)
	}
	return nil
}

// RemoveDislikeComment removes a dislike from a comment.
func (s *CommentStorage) RemoveDislikeComment(commentID int, username string) error {
	query := `DELETE FROM dislike WHERE commentId = $1 AND username = $2;`
	_, err := s.db.Exec(query, commentID, username)
	if err != nil {
		return fmt.Errorf("storage: remove like from comment: %w", err)
	}
	query = `UPDATE comment SET dislike = dislike - 1 WHERE id = $1;`
	_, err = s.db.Exec(query, commentID)
	if err != nil {
		return fmt.Errorf("storage: remove like from comment: %w", err)
	}
	return nil
}

// CommentHasLike checks if a comment has a like from a given user.
func (s *CommentStorage) CommentHasLike(commentID int, username string) error {
	var u, query string
	query = `SELECT username FROM like WHERE commentId = $1 AND username = $2;`
	err := s.db.QueryRow(query, commentID, username).Scan(&u)
	if err != nil {
		return fmt.Errorf("storage: comment has like: %w", err)
	}
	return nil
}

// CommentHasDislike checks if a comment has a dislike from a given user.
func (s *CommentStorage) CommentHasDislike(commentID int, username string) error {
	var u, query string
	query = `SELECT username FROM dislike WHERE commentId = $1 AND username = $2;`
	err := s.db.QueryRow(query, commentID, username).Scan(&u)
	if err != nil {
		return fmt.Errorf("storage: comment has like: %w", err)
	}
	return nil
}

// LikeComment adds a like to a comment.
func (s *CommentStorage) LikeComment(commentID int, username string) error {
	query := `INSERT INTO like(commentId, username) VALUES ($1, $2);`
	_, err := s.db.Exec(query, commentID, username)
	if err != nil {
		return fmt.Errorf("storage: like comment: %w", err)
	}
	query = `UPDATE comment SET like = like + 1  WHERE id = $1;`
	_, err = s.db.Exec(query, commentID)
	if err != nil {
		return fmt.Errorf("storage: like comment: %w", err)
	}
	return nil
}

// DislikeComment adds a dislike to a comment.
func (s *CommentStorage) DislikeComment(commentID int, username string) error {
	query := `INSERT INTO dislike(commentId, username) VALUES ($1, $2);`
	_, err := s.db.Exec(query, commentID, username)
	if err != nil {
		return fmt.Errorf("storage: like comment: %w", err)
	}
	query = `UPDATE comment SET dislike = dislike + 1  WHERE id = $1;`
	_, err = s.db.Exec(query, commentID)
	if err != nil {
		return fmt.Errorf("storage: like comment: %w", err)
	}
	return nil
}
