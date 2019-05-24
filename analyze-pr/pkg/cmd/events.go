package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/xeonx/timeago"
	"gopkg.in/AlecAivazis/survey.v1"
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

func EventsCommand(_ Context, jobURL string) error {
	content, err := fetchEventsArtifact(jobURL)
	if err != nil {
		return err
	}
	eventList := EventList{}
	if err := json.Unmarshal(content, &eventList); err != nil {
		return err
	}

	sort.Slice(eventList.Items, func(i, j int) bool { return eventList.Items[i].Timestamp.Before(eventList.Items[j].Timestamp) })

	englishFormat := timeago.English
	englishFormat.PastSuffix = " "
	w := tabwriter.NewWriter(os.Stdout, 60, 0, 0, ' ', tabwriter.DiscardEmptyColumns)

	black := color.New(color.FgHiBlack).SprintfFunc()
	white := color.New(color.FgWhite).SprintfFunc()
	red := color.New(color.FgRed).SprintfFunc()

	components := map[string]bool{}
	for _, item := range eventList.Items {
		if len(strings.TrimSpace(item.Source.Component)) == 0 {
			continue
		}
		components[item.Source.Component] = true
	}
	componentList := []string{}
	for component := range components {
		componentList = append(componentList, component)
	}
	sort.Strings(componentList)

	component := ""
	prompt := &survey.Select{
		Message:  "What component you want to see?",
		PageSize: 20,
		Options:  componentList,
	}
	if err := survey.AskOne(prompt, &component, nil); err != nil {
		return err
	}

	for _, item := range eventList.Items {
		if item.Source.Component != component {
			continue
		}
		humanTime := englishFormat.FormatReference(eventList.Items[0].Timestamp, item.Timestamp)

		message := item.Message
		if item.Type == "Warning" {
			message = red(message)
		}
		if _, err := fmt.Fprintf(w, "%s  %s\t%s\n", black(humanTime), white(item.Reason), message); err != nil {
			return err
		}
		if err := w.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func fetchEventsArtifact(u string) ([]byte, error) {
	originalURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(originalURL.Path, "/")
	newURL := "https://storage.googleapis.com/" + strings.Join(parts[2:], "/") + "/artifacts/e2e-aws/events.json"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	log.Printf("Fetching %s ...", newURL)
	response, err := client.Get(newURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, err
	}
	return ioutil.ReadAll(response.Body)
}
