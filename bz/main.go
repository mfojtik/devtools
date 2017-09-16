package main

import (
	"log"

	"github.com/mfojtik/devtools/bz/query"
	"github.com/mfojtik/devtools/bz/render"
	"github.com/mfojtik/devtools/bz/types"
)

func main() {
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

	render.Console(&bugList)
}
