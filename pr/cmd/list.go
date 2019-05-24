package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/mfojtik/devtools/pr/pkg/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List all GitHub pull requests for the user in the specified org",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		gh := client.GetGithubClientOrDie(ctx)

		prs, err := client.ListPullRequestsForUser(gh, "openshift", "mfojtik", false)
		if err != nil {
			log.Fatalf("Error fetching pull requests: %v", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 110, 0, 0, ' ', tabwriter.DiscardEmptyColumns)
		for _, pr := range prs {
			failingCount := 0
			successCount := 0
			pendingCount := 0
			for _, status := range pr.Statuses {
				switch status.GetState() {
				case "failure":
					failingCount++
				case "pending":
					pendingCount++
				case "success":
					successCount++
				}
			}
			status := fmt.Sprintf("f=%d,p=%d,s=%d", failingCount, pendingCount, successCount)
			title := pr.PullRequest.GetTitle()
			if len(title) > 80 {
				title = title[0:80] + " ..."
			}
			if _, err := fmt.Fprintf(w, "[ %8s ] [ %-15s ] %s\t | %s\n", pr.PullRequest.GetMergeableState(), status, pr.PullRequest.GetHTMLURL(), title); err != nil {
				log.Fatal(err)
			}
			if err := w.Flush(); err != nil {
				log.Fatal(err)
			}
		}
	},
}
