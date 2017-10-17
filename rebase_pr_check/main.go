package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func currentBranchName() string {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
	if err != nil {
		log.Fatalf("error running git: %q (%v)", out, err)
	}
	return strings.TrimSpace(string(out))
}

func isInGitLog(title string) bool {
	out, err := exec.Command("git", "log", "--all", "--grep", title).CombinedOutput()
	if err != nil {
		log.Fatalf("error running git: %q (%v)", out, err)
	}
	return len(strings.TrimSpace(string(out))) != 0
}

func main() {
	if len(os.Args) == 1 || len(os.Args[1]) == 0 {
		log.Fatalf("usage: %s FILE", os.Args[0])
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	prs := []int{}

	for scanner.Scan() {
		token := scanner.Text()
		parts := strings.Split(token, "UPSTREAM:")
		if len(parts) < 2 {
			log.Warnf("Failed to process: %s", token)
			continue
		}
		numParts := strings.Split(parts[1], ":")
		prNumber, err := strconv.Atoi(strings.TrimSpace(numParts[0]))
		if err != nil || prNumber == 0 {
			log.Warnf("skipping %s:%s ...", strings.TrimSpace(numParts[0]), numParts[1])
			continue
		}
		prs = append(prs, prNumber)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	currentBranch := currentBranchName()

	for _, prNumber := range prs {
		prLog := log.WithFields(log.Fields{
			"number": fmt.Sprintf("#%d", prNumber),
		})
		pr, _, err := client.PullRequests.Get(ctx, "kubernetes", "kubernetes", prNumber)
		if err != nil {
			prLog.Errorf("error fetching status: %v", err)
		}
		prLog = prLog.WithField("title", pr.GetTitle()).
			WithField("base", pr.Base.GetLabel())

		mergeStatus := "NOT MERGED"
		if *pr.Merged {
			mergeStatus = "MERGED"
			if isInGitLog(pr.GetTitle()) {
				mergeStatus += " (available in " + currentBranch + ")"
			}
		}
		prLog.Info(mergeStatus)
	}
}
