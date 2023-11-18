package service

import (
	"forum/internal/repository"
)

// Service is a struct that implements the Authorization, PostItem and Comment interfaces.
type Service struct {
	Authorization
	PostItem
	Comment
}

// NewService returns a new instance of Service.
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		PostItem:      NewPostService(repos.PostItem),
		Comment:       NewCommentService(repos.Comment),
	}
}
