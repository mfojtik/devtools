package client

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// shrug?
var archivedRepos = []string{
	"vagrant-openshift",
}

func GetGithubClientOrDie(ctx context.Context) *github.Client {
	if len(os.Getenv("GITHUB_TOKEN")) == 0 {
		panic("'GITHUB_TOKEN' environment variable is missing")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	return github.NewClient(oauth2.NewClient(ctx, ts))
}

func GetPullRequestStatus(gh *github.Client, pull *github.PullRequest) ([]*github.RepoStatus, error) {
	var statuses []*github.RepoStatus
	request, _ := gh.NewRequest("GET", pull.GetStatusesURL(), nil)
	if _, err := gh.Do(context.Background(), request, &statuses); err != nil {
		return nil, err
	}
	return statuses, nil
}

type PullRequestWithStatus struct {
	PullRequest *github.PullRequest
	Statuses    []*github.RepoStatus
}

func ListPullRequestsForUser(gh *github.Client, organization, username string, noStatus bool) ([]PullRequestWithStatus, error) {
	searchResult, _, err := gh.Search.Issues(context.Background(), "author:"+username+"+type:pr+state:open+org:"+organization, &github.SearchOptions{
		Sort:  "updated",
		Order: "desc",
	})
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	errChan := make(chan error)
	pullStatusChan := make(chan *github.PullRequest)

	var (
		result    []PullRequestWithStatus
		allErrors []error
	)

	// start collecting all errors
	go func() {
		for err := range errChan {
			allErrors = append(allErrors, err)
		}
	}()

	go func() {
		defer close(errChan)
		// wait for all pull request to have statuses fetched
		for pull := range pullStatusChan {
			var (
				states []*github.RepoStatus
				err    error
			)
			if !noStatus {
				states, err = GetPullRequestStatus(gh, pull)
				if err != nil {
					errChan <- err
				}
			}
			result = append(result, PullRequestWithStatus{PullRequest: pull, Statuses: states})
		}
	}()

	for _, i := range searchResult.Issues {
		repositoryParts := strings.Split(i.GetRepositoryURL(), "/")
		wg.Add(1)
		go func(repoName string, issue github.Issue) {
			defer wg.Done()
			for _, archived := range archivedRepos {
				if repoName == archived {
					return
				}
			}
			pull, _, err := gh.PullRequests.Get(context.Background(), organization, repoName, issue.GetNumber())
			if err != nil {
				errChan <- err
			} else {
				pullStatusChan <- pull
			}
		}(repositoryParts[len(repositoryParts)-1], i)
	}

	wg.Wait()
	close(pullStatusChan)

	<-errChan

	// report errors
	var aggregatedErrorMsg []string
	for _, err := range allErrors {
		aggregatedErrorMsg = append(aggregatedErrorMsg, err.Error())
	}
	if len(aggregatedErrorMsg) > 0 {
		return nil, fmt.Errorf(strings.Join(aggregatedErrorMsg, "\n"))
	}

	return result, nil
}
