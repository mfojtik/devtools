package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func newPeriodicScraper(duration time.Duration, stopCh <-chan struct{}) {
	immediateTicker := time.NewTicker(1 * time.Second)
	ticker := time.NewTicker(duration)
	tickChn := immediateTicker.C
	for {
		select {
		case <-tickChn:
			tickChn = ticker.C
			immediateTicker.Stop()
			newComponentBuilds := []Component{}
			for componentName, componentTestURL := range componentURLs {
				log.Printf("Updating data for %q ...", componentName)
				c := Component{Name: componentName, Builds: map[string]Build{}}
				builds, err := scrapeBuilds(componentTestURL)
				if err != nil {
					log.Printf("Error scraping new builds for %s: %v", componentName, err)
				} else {
					c.Builds = builds
				}
				newComponentBuilds = append(newComponentBuilds, c)
			}
			if err := syncData(newComponentBuilds); err != nil {
				log.Printf("Error storing data: %v", err)
			}
			log.Printf("Sleeping for %s ...", duration)
			// do stuff
		case <-stopCh:
			ticker.Stop()
			return
		}
	}
}

func syncData(newData []Component) error {
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		if err := os.Mkdir("data", os.ModePerm); err != nil {
			return err
		}
	}

	prevDataBytes, err := ioutil.ReadFile("data/data.json")
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	prevData := []Component{}

	if !os.IsNotExist(err) {
		if err := json.Unmarshal(prevDataBytes, &prevData); err != nil {
			return err
		}
	}

	currData := []Component{}

	for _, c := range prevData {
		newBuilds := map[string]Build{}
		for id, b := range c.Builds {
			// filter builds older than 3 weeks
			if time.Since(b.Timestamp) > time.Hour*24*21 {
				continue
			}
			newBuilds[id] = b
		}
		currData = append(currData, Component{Name: c.Name, Builds: newBuilds})
	}
	for _, c := range newData {
		newBuilds := map[string]Build{}
		for id, b := range c.Builds {
			// filter builds older than 3 weeks
			if time.Since(b.Timestamp) > time.Hour*24*21 {
				continue
			}
			newBuilds[id] = b
		}
		currData = append(currData, Component{Name: c.Name, Builds: newBuilds})
	}

	newDataBytes, err := json.Marshal(newData)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile("data/_data.json", newDataBytes, os.ModePerm); err != nil {
		return err
	}
	return os.Rename("data/_data.json", "data/data.json")
}

func scrapeBuilds(buildURL string) (map[string]Build, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(buildURL)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("unable to close body response: %v", err)
		}
	}()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	builds := map[string]Build{}
	doc.Find("#content a").Each(func(i int, s *goquery.Selection) {
		if s.Children().Length() != 2 {
			return
		}
		buildNumber := s.ChildrenFiltered(".build-number").Contents().Text()
		if len(buildNumber) == 0 {
			return
		}

		timestamp, exists := s.Find(".timestamp").Attr("data-epoch")
		if !exists || len(timestamp) == 0 {
			return
		}

		timestampInt, err := strconv.ParseInt(timestamp, 10, 32)
		if err != nil {
			return
		}

		var status Status
		switch {
		case s.Find(".build-timestamp").First().HasClass("build-failure"):
			status = StatusFailed
		case s.Find(".build-timestamp").First().HasClass("build-success"):
			status = StatusSuccess
		default:
			status = StatusPending

		}

		builds[buildNumber] = Build{
			Timestamp: time.Unix(timestampInt, 0),
			Status:    status,
		}
	})

	return builds, nil
}
