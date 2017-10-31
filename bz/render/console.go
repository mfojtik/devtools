package render

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/mfojtik/devtools/bz/types"
)

// Console renders the buglist to console.
func Console(bugList *types.Buglist) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Bug ID", "Assignee", "Status", "Summary"})
	table.SetFooter([]string{"", "", "Total", fmt.Sprintf("%d", len(bugList.Items))})
	table.SetBorder(false)
	table.SetColumnColor(
		tablewriter.Color(tablewriter.Bold, tablewriter.FgHiBlackColor),
		tablewriter.Color(tablewriter.Bold, tablewriter.FgHiRedColor),
		tablewriter.Color(tablewriter.Bold, tablewriter.FgHiBlackColor),
		tablewriter.Color(tablewriter.Bold, tablewriter.FgBlackColor),
	)
	for _, b := range bugList.Items {
		table.Append([]string{
			b.BugID,
			b.Assignee,
			b.Status,
			b.Summary,
		})
	}
	table.Render()

	return nil
}

// ConsoleDetails renders the buglist to console with details.
func ConsoleDetails(bugList *types.Buglist) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Bug ID", "Assignee", "Version", "Target", "Status", "Keywords", "Summary"})
	table.SetFooter([]string{"", "", "", "", "", "Total", fmt.Sprintf("%d", len(bugList.Items))})
	table.SetBorder(false)
	table.SetColumnColor(
		tablewriter.Color(tablewriter.Bold, tablewriter.FgHiBlackColor),
		tablewriter.Color(tablewriter.Bold, tablewriter.FgHiRedColor),
		tablewriter.Color(tablewriter.Bold, tablewriter.FgHiBlackColor),
		tablewriter.Color(tablewriter.Bold, tablewriter.FgHiBlackColor),
		tablewriter.Color(tablewriter.Bold, tablewriter.FgHiBlackColor),
		tablewriter.Color(tablewriter.Bold, tablewriter.FgHiBlackColor),
		tablewriter.Color(tablewriter.Bold, tablewriter.FgBlackColor),
	)
	for _, b := range bugList.Items {
		table.Append([]string{
			b.BugID,
			b.Assignee,
			b.Details.Version,
			b.Details.TargetRelease,
			b.Status,
			strings.Join(b.Details.Keywords, ","),
			b.Summary,
		})
	}
	table.Render()

	return nil
}
