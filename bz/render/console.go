package render

import (
	"fmt"
	"os"

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
