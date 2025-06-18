package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/google/go-github/v72/github"
)

func main() {
	ghToken := flag.String("token", "", "Github Token")
	if len(*ghToken) == 0 {
		fmt.Println("A GitHub token is required.")
		os.Exit(1)
	}

	prNumber := flag.Int("pr", 0, "Pull Request Number")
	if *prNumber == 0 {
		fmt.Println("A pull request number is required.")
		os.Exit(1)
	}

	repoOwner := flag.String("repoOwner", "", "Github Repo Owner")
	if len(*repoOwner) == 0 {
		fmt.Println("A repoOwner is required.")
		os.Exit(1)
	}

	repo := flag.String("repo", "", "Github Repo")
	if len(*repo) == 0 {
		fmt.Println("A repo is required.")
		os.Exit(1)
	}

	body := flag.String("body", "", "Body to send")
	if len(*body) == 0 {
		fmt.Println("A body is required.")
		os.Exit(1)
	}

	commentId := flag.Int64("comment-id", 0, "Comment ID")
	bodyIncludes := flag.String("body-includes", "", "Comment ID")

	client := github.NewClient(nil).WithAuthToken(*ghToken)

	if *commentId == 0 && len(*bodyIncludes) > 0 {
		*commentId = findComment(client, *repoOwner, *repo, *prNumber, *bodyIncludes)
	}

	if *commentId == 0 {
		createComment(client, *repoOwner, *repo, *prNumber, *body)
	} else {
		updateComment(client, *repoOwner, *repo, *commentId, *body)
	}
}

func createComment(client *github.Client, repoOwner string, repo string, prNumber int, comment string) {
	_, _, err := client.Issues.CreateComment(context.Background(), repoOwner, repo, prNumber, &github.IssueComment{Body: github.Ptr(comment)})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Comment created...")
}

func updateComment(client *github.Client, repoOwner string, repo string, commentId int64, comment string) {
	_, _, err := client.Issues.EditComment(context.Background(), repoOwner, repo, commentId, &github.IssueComment{Body: github.Ptr(comment)})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Comment updated...")
}

func findComment(client *github.Client, repoOwner string, repo string, prNumber int, bodyIncludes string) int64 {
	comments, _, err := client.Issues.ListComments(context.Background(), repoOwner, repo, prNumber, &github.IssueListCommentsOptions{})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	matchingComments := slices.Collect(
		func(yield func(*github.IssueComment) bool) {
			for _, comment := range comments {
				if strings.Contains(comment.GetBody(), bodyIncludes) {
					if !yield(comment) {
						return
					}
				}
			}
		},
	)

	if len(matchingComments) > 0 {
		fmt.Println("Comment found...")

		return matchingComments[0].GetID()
	} else {
		fmt.Println("Comment not found...")

		return 0
	}
}
