package cmd

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/go-github/github"
	"github.com/xeonx/timeago"
	"gopkg.in/AlecAivazis/survey.v1"
)

func StatusCommand(ctx Context) error {
	statuses := []*github.RepoStatus{}
	request, err := ctx.Client.NewRequest("GET", ctx.PullRequest.GetStatusesURL(), nil)
	if err != nil {
		return err
	}
	if _, err := ctx.Client.Do(context.Background(), request, &statuses); err != nil {
		return err
	}

	sort.Slice(statuses, func(i, j int) bool { return statuses[i].UpdatedAt.After(*statuses[j].UpdatedAt) })

	jobList := map[string]github.RepoStatus{}
	finishedJobKeys := []string{}

	for _, s := range statuses {
		// skip pending checks as there are no artifacts for them yet
		if s.State == nil || s.GetState() == "pending" {
			continue
		}
		// skip tide checks ...
		if s.GetContext() == "tide" {
			continue
		}
		k := fmt.Sprintf("%s [%s %s]", s.GetContext(), s.GetState(), timeago.English.Format(*s.UpdatedAt))
		jobList[k] = *s
		finishedJobKeys = append(finishedJobKeys, k)
	}

	prompt := &survey.Select{
		Message:  "Pick a CI job:",
		PageSize: 20,
		Options:  finishedJobKeys,
	}
	job := ""
	if err := survey.AskOne(prompt, &job, nil); err != nil {
		return err
	}

	action := ""
	prompt = &survey.Select{
		Message: "What information you want to see?",
		Options: []string{
			"Events",
			"Logs",
		},
	}
	if err := survey.AskOne(prompt, &action, nil); err != nil {
		return err
	}

	switch action {
	case "Events":
		if err := EventsCommand(ctx, *jobList[job].TargetURL); err != nil {
			return err
		}
	case "Logs":
		if err := LogsCommand(ctx, *jobList[job].TargetURL); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	return nil
}
