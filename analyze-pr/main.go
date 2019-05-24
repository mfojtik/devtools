package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/mfojtik/devtools/analyze-pr/pkg/cmd"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if len(token) == 0 {
		log.Fatalf("Please set the 'GITHUB_TOKEN' environment variable.")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	if len(os.Args) == 1 {
		log.Fatalf("Usage: %s [PULL_REQUEST_URL]", os.Args[0])
	}

	prURL, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid pull request URL")
	}

	parts := strings.Split(prURL.Path, "/")

	prNum, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		log.Fatalf("Wrong pull number: %v", err)
	}
	pr, _, err := client.PullRequests.Get(ctx, parts[len(parts)-4], parts[len(parts)-3], prNum)
	if err != nil {
		log.Fatalf("Unable to get pull request: %v", err)
	}

	commandCtx := cmd.Context{
		PullRequest: *pr,
		Client:      client,
	}

	for {
		command, err := cmd.ReadCommand(commandCtx)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}

		switch command {
		case "status":
			if err := cmd.StatusCommand(commandCtx); err != nil {
				log.Printf("error: %v", err)
			}
		case "":
			continue
		default:
			log.Printf("unknown command: %q", command)
		}
	}

}
