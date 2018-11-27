package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/xeonx/timeago"
)

type EventList struct {
	Items []Event `json:"items"`
}

type Event struct {
	Timestamp time.Time `json:"firstTimestamp"`
	Reason    string    `json:"reason"`
	Message   string    `json:"message"`
	Source    Source    `json:"source"`
	Type      string    `json:"type"`
}

type Source struct {
	Component string `json:"component"`
}

func getEvents(u string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Get(u)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Fatalf("get %s failed with: %d", u, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("usage: %s <events.json> <component>", os.Args[0])
	}
	content, err := getEvents(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}
	/*
		content, err := ioutil.ReadFile(os.Args[1])

	*/
	eventList := EventList{}
	if err := json.Unmarshal(content, &eventList); err != nil {
		log.Fatal(err.Error())
	}

	sort.Slice(eventList.Items, func(i, j int) bool { return eventList.Items[i].Timestamp.Before(eventList.Items[j].Timestamp) })

	englishFormat := timeago.English
	englishFormat.PastSuffix = " "
	w := tabwriter.NewWriter(os.Stdout, 60, 0, 0, ' ', tabwriter.DiscardEmptyColumns)

	black := color.New(color.FgHiBlack).SprintfFunc()
	white := color.New(color.FgWhite).SprintfFunc()
	red := color.New(color.FgRed).SprintfFunc()

	for _, item := range eventList.Items {
		if item.Source.Component != os.Args[2] {
			continue
		}
		humanTime := englishFormat.FormatReference(eventList.Items[0].Timestamp, item.Timestamp)

		message := item.Message
		if item.Type == "Warning" {
			message = red(message)
		}
		fmt.Fprintf(w, "%s  %s\t%s\n", black(humanTime), white(item.Reason), message)
		w.Flush()
	}
}
