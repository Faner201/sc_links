package github

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

type Client struct{}

func NewClien() *Client {
	return &Client{}
}

func (c *Client) ExchangeCodeToAccessKey(ctx context.Context, clientID, clientSecret, code string) (string, error) {
	type exchangeCodeRequest struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
	}

	var respJSON struct {
		AccessToken string `json:"access_token"`
	}

	req := exchangeCodeRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Code:         code,
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://github.com/login/oauth/access_token",
		bytes.NewReader(reqJSON),
	)
	if err != nil {
		return "", err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil {
		return "", err
	}

	return respJSON.AccessToken, nil
}

func (c *Client) IsMember(ctx context.Context, accessKey, org, user string) (bool, error) {
	githubClient := getGithubClientWithAccessKey(ctx, accessKey)
	isMember, _, err := githubClient.Organizations.IsMember(ctx, org, user)
	if err != nil {
		return false, err
	}
	return isMember, err
}

func (c *Client) GetUser(ctx context.Context, accessKey, user string) (*github.User, error) {
	githubClient := getGithubClientWithAccessKey(ctx, accessKey)
	ghUser, _, err := githubClient.Users.Get(ctx, user)
	if err != nil {
		return nil, err
	}

	return ghUser, nil
}

func getGithubClientWithAccessKey(ctx context.Context, accessKey string) *github.Client {
	var (
		ts = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessKey})
		tc = oauth2.NewClient(ctx, ts)
	)

	return github.NewClient(tc)
}
