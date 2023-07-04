package user

import (
	"context"
	"fmt"

	"github.com/Faner201/sc_links/internal/db"
	errors "github.com/Faner201/sc_links/internal/error"
	"github.com/Faner201/sc_links/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mgo struct {
	db *mongo.Database
}

func NewMongoDB(client *mongo.Database) *mgo {
	return &mgo{db: client}
}

func (m *mgo) col() *mongo.Collection {
	return m.db.Collection("users")
}

func (m *mgo) update(ctx context.Context, user model.User, upsert bool) error {
	var (
		query       = bson.M{"_id": user.GitHubLogin}
		replacement = mgoUserFromModel(user)
		opts        = &options.ReplaceOptions{Upsert: &upsert}
	)

	if _, err := m.col().ReplaceOne(ctx, query, replacement, opts); err != nil {
		return err
	}

	return nil
}

func (m *mgo) CreateOrUpdate(ctx context.Context, user model.User) (*model.User, error) {
	const operation = "users.mgo.CreateOrUpdate"

	if err := m.update(ctx, user, true); err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &user, nil
}

func (m *mgo) Update(ctx context.Context, user model.User) error {
	const operation = "users.mgo.Update"

	if err := m.update(ctx, user, false); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

func (m *mgo) Get(ctx context.Context, ghLogin string) (*model.User, error) {
	const operation = "users.mgo.Get"

	var user db.MgoUser
	if err := m.col().FindOne(ctx, bson.M{"_id": ghLogin}).Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%s: %w", operation, errors.ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return modelUserFromMgo(user), nil
}

func (m *mgo) GetByGithubLogin(ctx context.Context, ghLogin string) (*model.User, error) {
	const operation = "usre.mgo.GetByGithubLogin"

	var user db.MgoUser
	if err := m.col().FindOne(ctx, bson.M{"_id": ghLogin}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%s: %w", operation, err)
		}
	}

	return modelUserFromMgo(user), nil
}

func (m *mgo) Deactivate(ctx context.Context, ghLogin string) error {
	const operation = "users.mgo.Deactivate"

	user, err := m.Get(ctx, ghLogin)

	if err != nil {
		return fmt.Errorf("%s: %w", operation, errors.ErrNotFound)
	}

	user.IsActive = false

	if err := m.update(ctx, *user, false); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

func mgoUserFromModel(m model.User) db.MgoUser {
	return db.MgoUser{
		IsActive:        m.IsActive,
		GitHubLogin:     m.GitHubLogin,
		GitHubAccessKey: m.GitHubAccessKey,
		CreatedAt:       m.CreatedAt,
	}
}

func modelUserFromMgo(m db.MgoUser) *model.User {
	return &model.User{
		IsActive:        m.IsActive,
		GitHubLogin:     m.GitHubLogin,
		GitHubAccessKey: m.GitHubAccessKey,
		CreatedAt:       m.CreatedAt,
	}
}
