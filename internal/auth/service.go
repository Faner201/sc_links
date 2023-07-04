package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/Faner201/sc_links/internal/config"
	errors "github.com/Faner201/sc_links/internal/error"
	"github.com/Faner201/sc_links/internal/model"
	"github.com/Faner201/sc_links/internal/storage"
	"github.com/google/go-github/v48/github"
)

//go:generate moq -out mock_github_client.go -pkg auth . GithubClient
type GithubClient interface {
	ExchangeCodeToAccessKey(ctx context.Context, clientID, clientSecret, code string) (string, error)
	IsMember(ctx context.Context, accessKey, org, user string) (bool, error)
	GetUser(ctx context.Context, accessKey, user string) (*github.User, error)
}

type Service struct {
	github  GithubClient
	storage storage.StorageGithub

	ghClientID     string
	ghClientSecret string
}

func NewService(githubClient GithubClient, storage storage.StorageGithub, ghClientID, ghClientSecret string) *Service {
	return &Service{
		github:         githubClient,
		storage:        storage,
		ghClientID:     ghClientID,
		ghClientSecret: ghClientSecret,
	}
}

func (s *Service) GithubAuthLink() string {
	return fmt.Sprintf("https://github.com/login/oauth/authorize?scopes=user,read:org&client_id=%s", s.ghClientID)
}

func (s *Service) GithubAuthCallback(ctx context.Context, sessionCode string) (*model.User, string, error) {
	accessKey, err := s.github.ExchangeCodeToAccessKey(ctx, s.ghClientID, s.ghClientSecret, sessionCode)
	if err != nil {
		return nil, "", err
	}

	ghUser, err := s.github.GetUser(ctx, accessKey, "")
	if err != nil {
		return nil, "", err
	}

	user, err := s.RegisterUser(ctx, ghUser, accessKey)
	if err != nil {
		return nil, "", err
	}

	jwt, err := MakeJWT(*user)
	if err != nil {
		log.Printf("failed to make jwt: %v", err)
		return nil, "", err
	}

	return user, jwt, nil

}

func (s *Service) RegisterUser(ctx context.Context, ghUser *github.User, accessKey string) (*model.User, error) {
	isMember, err := s.github.IsMember(ctx, accessKey, config.Get().Auth.AllowedGithubOrg, ghUser.GetLogin())
	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, fmt.Errorf("%w %q", errors.ErrUserIsNotMember, config.Get().Auth.AllowedGithubOrg)
	}

	user := model.User{
		GitHubLogin:     ghUser.GetLogin(),
		IsActive:        true,
		GitHubAccessKey: accessKey,
	}

	return s.storage.CreateOrUpdate(ctx, user)
}
