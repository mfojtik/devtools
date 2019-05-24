package cmd

import (
	"github.com/google/go-github/github"
)

type Context struct {
	Client      *github.Client
	PullRequest github.PullRequest
}
