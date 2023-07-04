package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Faner201/sc_links/internal/auth"
	er "github.com/Faner201/sc_links/internal/error"
	"github.com/Faner201/sc_links/internal/model"
	"github.com/Faner201/sc_links/internal/storage/user"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/go-github/v48/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGithubAuthLink(t *testing.T) {
	t.Run("returns github auth link with client id", func(t *testing.T) {
		svc := auth.NewService(nil, nil, "client-id", "")
		assert.Equal(
			t,
			"https://github.com/login/oauth/authorize?scopes=user,read:org&client_id=client-id",
			svc.GithubAuthLink(),
		)
	})
}

func TestGithubAuthCallback(t *testing.T) {
	t.Run("returns user model and JWT", func(t *testing.T) {
		var (
			ghClient = &auth.GithubClientMock{
				ExchangeCodeToAccessKeyFunc: func(ctx context.Context, clientID, clientSecret, code string) (string, error) {
					return "access-key", nil
				},
				GetUserFunc: func(ctx context.Context, accessKey, user string) (*github.User, error) {
					return &github.User{
						Login: github.String(gofakeit.Username()),
					}, nil
				},
				IsMemberFunc: func(ctx context.Context, accessKey, org, user string) (bool, error) {
					return true, nil
				},
			}
			userStorage = user.NewInMemory()
			svc         = auth.NewService(ghClient, userStorage, "", "")
		)

		user, token, err := svc.GithubAuthCallback(context.Background(), gofakeit.Numerify("code-###"))
		require.NoError(t, err)
		assert.True(t, user.IsActive)
		assert.NotEmpty(t, token)
	})

	t.Run("returns error", func(t *testing.T) {
		t.Run("when exchanging code to access key fails", func(t *testing.T) {
			var (
				ghClient = &auth.GithubClientMock{
					ExchangeCodeToAccessKeyFunc: func(ctx context.Context, clientID, clientSecret, code string) (string, error) {
						return "", errors.New("exchange code to access key error")
					},
				}
				userStorage = user.NewInMemory()
				svc         = auth.NewService(ghClient, userStorage, "", "")
			)

			user, token, err := svc.GithubAuthCallback(context.Background(), gofakeit.Numerify("code-###"))
			assert.Error(t, err)
			assert.Nil(t, user)
			assert.Empty(t, token)
		})

		t.Run("when getting user from github fails", func(t *testing.T) {
			var (
				ghClient = &auth.GithubClientMock{
					ExchangeCodeToAccessKeyFunc: func(ctx context.Context, clientID, clientSecret, code string) (string, error) {
						return "access-key", nil
					},
					GetUserFunc: func(ctx context.Context, accessKey, user string) (*github.User, error) {
						return &github.User{
							Login: github.String(gofakeit.Username()),
						}, nil
					},
					IsMemberFunc: func(ctx context.Context, accessKey, org, user string) (bool, error) {
						return false, errors.New("is member error")
					},
				}
				userStorage = user.NewInMemory()
				svc         = auth.NewService(ghClient, userStorage, "", "")
			)

			user, token, err := svc.GithubAuthCallback(context.Background(), gofakeit.Numerify("code-###"))
			assert.Error(t, err)
			assert.Nil(t, user)
			assert.Empty(t, token)
		})

		t.Run("when registering user fails", func(t *testing.T) {
			var (
				ghClient = &auth.GithubClientMock{
					ExchangeCodeToAccessKeyFunc: func(ctx context.Context, clientID, clientSecret, code string) (string, error) {
						return "access-key", nil
					},
					GetUserFunc: func(ctx context.Context, accessKey, user string) (*github.User, error) {
						return nil, errors.New("get user error")
					},
				}
				userStorage = user.NewInMemory()
				svc         = auth.NewService(ghClient, userStorage, "", "")
			)

			user, token, err := svc.GithubAuthCallback(context.Background(), gofakeit.Numerify("code-###"))
			assert.Error(t, err)
			assert.Nil(t, user)
			assert.Empty(t, token)
		})
	})
}

func TestRegisterUser(t *testing.T) {
	t.Run("return user model", func(t *testing.T) {
		t.Run("when user is member of the organization", func(t *testing.T) {
			var (
				ghClient = &auth.GithubClientMock{
					IsMemberFunc: func(ctx context.Context, accessKey, org, user string) (bool, error) {
						return true, nil
					},
				}
				userStorage = user.NewInMemory()
				svc         = auth.NewService(ghClient, userStorage, "", "")
				ghUser      = &github.User{Login: github.String(gofakeit.Username())}
			)

			user, err := svc.RegisterUser(context.Background(), ghUser, "")
			require.NoError(t, err)
			assert.Equal(t, *ghUser.Login, user.GitHubLogin)
			assert.True(t, user.IsActive)
		})

		t.Run("even if user already exists", func(t *testing.T) {
			var (
				ghClient = &auth.GithubClientMock{
					IsMemberFunc: func(ctx context.Context, accessKey, org, user string) (bool, error) {
						return true, nil
					},
				}
				userStorage  = user.NewInMemory()
				svc          = auth.NewService(ghClient, userStorage, "", "")
				ghUser       = &github.User{Login: github.String(gofakeit.Username())}
				existingUser = model.User{
					IsActive:    true,
					GitHubLogin: ghUser.GetLogin(),
				}
			)

			_, err := userStorage.CreateOrUpdate(context.Background(), existingUser)
			require.NoError(t, err)

			user, err := svc.RegisterUser(context.Background(), ghUser, "")
			require.NoError(t, err)
			assert.Equal(t, *ghUser.Login, user.GitHubLogin)
			assert.True(t, user.IsActive)
		})
	})

	t.Run("return error", func(t *testing.T) {
		t.Run("when user is not a member of the organization", func(t *testing.T) {
			var (
				ghClient = &auth.GithubClientMock{
					IsMemberFunc: func(ctx context.Context, accessKey, org, user string) (bool, error) {
						return false, nil
					},
				}
				userStorage = user.NewInMemory()
				svc         = auth.NewService(ghClient, userStorage, "", "")
				ghUser      = &github.User{Login: github.String(gofakeit.Username())}
			)

			user, err := svc.RegisterUser(context.Background(), ghUser, "")
			assert.ErrorIs(t, err, er.ErrUserIsNotMember)
			assert.Nil(t, user)
		})

		t.Run("when github client returns error", func(t *testing.T) {
			var (
				ghClient = &auth.GithubClientMock{
					IsMemberFunc: func(ctx context.Context, accessKey, org, user string) (bool, error) {
						return false, nil
					},
				}
				userStorage = user.NewInMemory()
				svc         = auth.NewService(ghClient, userStorage, "", "")
				ghUser      = &github.User{Login: github.String(gofakeit.Username())}
			)

			user, err := svc.RegisterUser(context.Background(), ghUser, "")
			assert.Error(t, err)
			assert.Nil(t, user)
		})
	})
}
