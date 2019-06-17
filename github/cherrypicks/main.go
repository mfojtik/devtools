package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var repositories = []string{
	"origin",
	"cluster-kube-apiserver-operator",
	"cluster-kube-controller-manager-operator",
	"cluster-openshift-apiserver-operator",
	"cluster-kube-scheduler-operator",
	"cluster-openshift-controller-manager-operator",
	"descheduler-operator",
	"cluster-authentication-operator",
	"cluster-machine-approver",
	"service-ca-operator",
	"oauth-proxy",
	"service-serving-cert-signer",
}

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	for _, repo := range repositories {
		prs, _, err := client.PullRequests.List(ctx, "openshift", repo, &github.PullRequestListOptions{
			State: "open",
			Base:  "release-4.1",
			Sort:  "updated",
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		})
		if err != nil {
			panic(err)
		}
		for _, pr := range prs {
			labels := []string{}
			for _, label := range pr.Labels {
				labels = append(labels, label.GetName())
			}
			fmt.Printf("[%s/%s] #%d: %s (%s)\n", repo, pr.GetBase().GetRef(), pr.GetNumber(), pr.GetTitle(), strings.Join(labels, ","))
		}
	}
}
