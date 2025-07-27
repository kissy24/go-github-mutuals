package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v50/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

// Create GitHub Contents.
func CreateGHContents() (context.Context, *github.Client) {
	godotenv.Load()
	token := os.Getenv("GITHUB_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return ctx, client
}

// getAllUsers retrieves the complete list of followers or following for the specified user.
func getAllUsers(ctx context.Context, client *github.Client, userName string, listFunc func(context.Context, string, *github.ListOptions) ([]*github.User, *github.Response, error)) ([]*github.User, error) {
	var allUsers []*github.User
	opt := &github.ListOptions{Page: 1, PerPage: 100} // Increase the number of items per page
	for {
		users, resp, err := listFunc(ctx, userName, opt)
		if err != nil {
			return nil, err
		}
		allUsers = append(allUsers, users...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allUsers, nil
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
	followers, err := getAllUsers(ctx, client, userName, client.Users.ListFollowers)
	if err != nil {
		fmt.Printf("An error occurred while acquiring followers: %v\n", err)
		os.Exit(1)
	}
	followings, err := getAllUsers(ctx, client, userName, client.Users.ListFollowing)
	if err != nil {
		fmt.Printf("An error occurred while acquiring the following: %v\n", err)
		os.Exit(1)
	}
	fs := FollowStatus{followings, followers}
	notFollowingBack := fs.CompareUsers()
	fmt.Println("Users not following back:")
	for _, user := range notFollowingBack {
		fmt.Println(*user.Login)
	}
}
