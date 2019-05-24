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
	rootCmd.AddCommand(statusCommand)
}

var statusCommand = &cobra.Command{
	Use:   "status",
	Short: "Print pull request status",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		gh := client.GetGithubClientOrDie(ctx)

		prs, err := client.ListPullRequestsForUser(gh, "openshift", "mfojtik", false)
		if err != nil {
			log.Fatalf("Error fetching pull requests: %v", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 60, 0, 0, ' ', tabwriter.DiscardEmptyColumns)
		for _, pr := range prs {
			if _, err := fmt.Fprintf(w, "%s \t (%s)\n", pr.PullRequest.GetTitle(), pr.PullRequest.GetHTMLURL()); err != nil {
				log.Fatal(err)
			}
			for _, status := range pr.Statuses {
				// skip success
				if status.GetContext() == "success" {
					continue
				}
				if _, err := fmt.Fprintf(w, "  %s \t %s (%s)\n", status.GetContext(), status.GetState(), status.GetTargetURL()); err != nil {
					log.Fatal(err)
				}
			}
			if err := w.Flush(); err != nil {
				log.Fatal(err)
			}
		}
	},
}
