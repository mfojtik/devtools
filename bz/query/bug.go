package query

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/mfojtik/devtools/bz/types"
)

// FetchBugDetails retrieve the details about bug from remote
func FetchBugDetails(bugID string) (*types.BugDetails, error) {
	bugURL, _ := url.Parse("https://bugzilla.redhat.com/show_bug.cgi")
	params := url.Values{}
	params.Set("ctype", "xml")
	params.Set("id", bugID)
	bugURL.RawQuery = params.Encode()

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	resp, err := client.Get(bugURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := types.Bugzilla{}
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result.Bug, nil
}

// ProcessBuglistForDetails processes bug list and will fill it with details.
func ProcessBuglistForDetails(bugs *types.Buglist) error {
	workerChan := make(chan types.Bug)
	errors := []error{}
	worker := func(b types.Bug) {
		updated := b
		defer func() {
			workerChan <- updated
		}()
		details, err := FetchBugDetails(updated.BugID)
		if err != nil {
			errors = append(errors, err)
			return
		}
		updated.Details = details
	}

	for i := range bugs.Items {
		go worker(bugs.Items[i])
	}

	result := []types.Bug{}
	for {
		result = append(result, <-workerChan)
		if len(result) == len(bugs.Items) {
			break
		}
	}

	if len(errors) > 0 {
		result := []string{}
		for _, e := range errors {
			result = append(result, e.Error())
		}
		return fmt.Errorf("%s", strings.Join(result, ","))
	}
	bugs.Items = result

	return nil
}
