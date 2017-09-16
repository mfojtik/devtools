package query

import (
	"crypto/tls"
	"encoding/csv"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mfojtik/devtools/bz/types"
)

var (
	baseBzURL, _         = url.Parse("https://bugzilla.redhat.com/buglist.cgi")
	baseBzClassification = "Red Hat"
)

// BugzillaQueryBuilder contains configuration for RH bugzilla query builder.
type BugzillaQueryBuilder struct {
	severities []string
	statuses   []string
	components []string

	url    *url.URL
	values url.Values
	result *types.Buglist
}

// NewBugzillaQuery constructs and perform queries against RH bugzilla.
func NewBugzillaQuery() *BugzillaQueryBuilder {
	return &BugzillaQueryBuilder{
		url:    baseBzURL,
		values: url.Values{},
	}
}

// Complete completes the URL building part and returns the final query URL.
func (b *BugzillaQueryBuilder) Complete() *url.URL {
	b.values.Set("classification", baseBzClassification)
	b.values.Set("query_format", "advanced")
	b.values.Set("ctype", "csv")
	b.values.Set("f1", "version")
	b.values.Set("o1", "notregexp")
	b.values.Set("human", "1")
	b.values.Set("v1", `^2\.`)
	b.url.RawQuery = b.values.Encode()
	return b.url
}

func (b *BugzillaQueryBuilder) AddKeyword(name, keywordType string) *BugzillaQueryBuilder {
	b.values.Add("keywords", name)
	b.values.Add("keywords_type", keywordType)
	return b
}

func (b *BugzillaQueryBuilder) AddProducts(p ...string) *BugzillaQueryBuilder {
	for _, i := range p {
		b.values.Add("product", i)
	}
	return b
}

func (b *BugzillaQueryBuilder) AddSeverities(s ...string) *BugzillaQueryBuilder {
	for _, i := range s {
		b.values.Add("bug_severity", i)
	}
	return b
}

func (b *BugzillaQueryBuilder) AddStatuses(s ...string) *BugzillaQueryBuilder {
	for _, i := range s {
		b.values.Add("bug_status", i)
	}
	return b
}

func (b *BugzillaQueryBuilder) AddComponents(s ...string) *BugzillaQueryBuilder {
	for _, i := range s {
		b.values.Add("component", i)
	}
	return b
}

// Into sets the return object to fill.
func (b *BugzillaQueryBuilder) Into(list *types.Buglist) *BugzillaQueryBuilder {
	b.result = list
	return b
}

// Do performs the HTTP requests and parsing of CSV.
func (b *BugzillaQueryBuilder) Do() error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(b.Complete().String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	r := csv.NewReader(resp.Body)
	skipFirst := true
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if skipFirst {
			skipFirst = false
			continue
		}
		parsedTime, err := time.Parse("2006-01-02 15:04:05", record[7])
		b.result.Items = append(b.result.Items, types.Bug{
			BugID:      record[0],
			Product:    record[1],
			Component:  record[2],
			Assignee:   record[3],
			Status:     record[4],
			Resolution: record[5],
			Summary:    record[6],
			Changed:    parsedTime,
		})
	}
	return nil
}
