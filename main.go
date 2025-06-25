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
	ghToken := flag.String("token", "", "Github token")
	prNumber := flag.Int("pr", 0, "Pull request number")
	owner := flag.String("owner", "", "GitHub repository owner")
	repo := flag.String("repo", "", "GitHub repository name")
	body := flag.String("body", "", "Body to send")
	commentId := flag.Int64("comment-id", 0, "Comment ID")
	bodyIncludes := flag.String("body-includes", "", "Comment ID")

	flag.Parse()

	if len(*ghToken) == 0 {
		fmt.Println("A GitHub token is required.")
		os.Exit(1)
	}

	if *prNumber == 0 {
		fmt.Println("A pull request number is required.")
		os.Exit(1)
	}

	if len(*owner) == 0 {
		fmt.Println("A repository owner is required.")
		os.Exit(1)
	}

	if len(*repo) == 0 {
		fmt.Println("A repo is required.")
		os.Exit(1)
	}

	if len(*body) == 0 {
		fmt.Println("A body is required.")
		os.Exit(1)
	}

	client := github.NewClient(nil).WithAuthToken(*ghToken)

	if *commentId == 0 && len(*bodyIncludes) > 0 {
		*commentId = findComment(client, *owner, *repo, *prNumber, *bodyIncludes)
	}

	if *commentId == 0 {
		createComment(client, *owner, *repo, *prNumber, *body)
	} else {
		updateComment(client, *owner, *repo, *commentId, *body)
	}
}

func createComment(client *github.Client, owner string, repo string, prNumber int, comment string) {
	_, _, err := client.Issues.CreateComment(context.Background(), owner, repo, prNumber, &github.IssueComment{Body: github.Ptr(comment)})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Comment created...")
}

func updateComment(client *github.Client, owner string, repo string, commentId int64, comment string) {
	_, _, err := client.Issues.EditComment(context.Background(), owner, repo, commentId, &github.IssueComment{Body: github.Ptr(comment)})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Comment updated...")
}

func findComment(client *github.Client, owner string, repo string, prNumber int, bodyIncludes string) int64 {
	comments, _, err := client.Issues.ListComments(context.Background(), owner, repo, prNumber, &github.IssueListCommentsOptions{})

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
