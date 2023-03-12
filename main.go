package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

// Create GitHub Contents.
func CreateGHContents() (context.Context, *github.Client) {
	token := os.Getenv("GITHUB_TOKEN")
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
	userName := "kissy24" // Change your user name.
	ctx, client := CreateGHContents()
	followers, resp, err := client.Users.ListFollowers(ctx, userName, nil)
	if err != nil {
		fmt.Println(resp.Body)
		fmt.Println(err)
	}
	followings, resp, err := client.Users.ListFollowing(ctx, userName, nil)
	if err != nil {
		fmt.Println(resp.Body)
		fmt.Println(err)
	}
	fs := FollowStatus{followings, followers}
	notFollowingBack := fs.CompareUsers()
	fmt.Println("Users not following back:")
	for _, user := range notFollowingBack {
		fmt.Println(*user.Login)
	}
}
