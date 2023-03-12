package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Create GitHub Contents.
func CreateGHContents(g_token string) (context.Context, *github.Client) {
	token := os.Getenv(g_token)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return ctx, client
}

type FollowStatus struct {
	followings []*github.User
	followers  []*github.User
}

// Get a list of users who do not follow each other
func (fs FollowStatus) CompareUsers() []*github.User {
	var notFollowingBack []*github.User
	for _, user := range fs.followings {
		found := false
		for _, follower := range fs.followers {
			if *user.Login == *follower.Login {
				found = true
				break
			}
		}
		if !found {
			notFollowingBack = append(notFollowingBack, user)
		}
	}
	return notFollowingBack
}

func main() {
	token := "" // Set GitHub Token
	ctx, client := CreateGHContents(token)
	followers, _, err := client.Users.ListFollowers(ctx, "username", nil)
	if err != nil {
		fmt.Println(err)
	}
	followings, _, err := client.Users.ListFollowing(ctx, "username", nil)
	if err != nil {
		fmt.Println(err)
	}
	fs := FollowStatus{followings, followers}
	notFollowingBack := fs.CompareUsers()
	fmt.Println("Users not following back:")
	for _, user := range notFollowingBack {
		fmt.Println(*user.Login)
	}
}
