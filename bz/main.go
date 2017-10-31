package main

import (
	"flag"
	"log"

	"github.com/mfojtik/devtools/bz/query"
	"github.com/mfojtik/devtools/bz/render"
	"github.com/mfojtik/devtools/bz/types"
)

var detailsFlag = flag.Bool("details", false, "Show more details")

func main() {
	flag.Parse()
	var bugList types.Buglist
	b := query.NewBugzillaQuery()

	if err := b.AddSeverities("unspecified", "urgent", "high", "medium").
		AddStatuses("NEW", "ASSIGNED", "POST", "ON_DEV").
		AddKeyword("UpcomingRelease", "nowords").
		AddComponents("Deployments", "Master", "etcd").
		AddProducts("OpenShift Container Platform").
		AddProducts("OpenShift Online").
		AddProducts("OpenShift Origin").
		Into(&bugList).Do(); err != nil {
		log.Fatalf("Unable to fetch bugs from %s: %v", b.Complete().String(), err)
	}

	if *detailsFlag == true {
		if err := query.ProcessBuglistForDetails(&bugList); err != nil {
			log.Fatalf("Unable to process bug details: %v", err)
		}
		render.ConsoleDetails(&bugList)
		return
	}
	render.Console(&bugList)
}
