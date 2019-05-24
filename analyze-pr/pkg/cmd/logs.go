package cmd

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"gopkg.in/AlecAivazis/survey.v1"
)

func parseJobURL(pullRequestLink string) string {
	newURL := strings.Replace(pullRequestLink, "openshift-gce-devel.appspot.com/build/", "gcsweb-ci.svc.ci.openshift.org/gcs/", 1)
	return newURL + "/artifacts/e2e-aws/pods/"
}

func parseArtifactsURLs(responseBody io.Reader) []string {
	tokenizer := html.NewTokenizer(responseBody)
	var urls []string
	for {
		tt := tokenizer.Next()
		switch {
		case tt == html.ErrorToken:
			return urls
		case tt == html.StartTagToken:
			t := tokenizer.Token()
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}
			for _, attr := range t.Attr {
				if attr.Key == "href" && strings.Contains(attr.Val, "/artifacts/") {
					urls = append(urls, attr.Val)
					break
				}
			}
		}
	}
}

func fetchArtifacts(urls []string) map[string][]byte {
	var wg sync.WaitGroup
	result := map[string][]byte{}
	for _, u := range urls {
		artifactURL, err := url.Parse(u)
		if err != nil {
			log.Printf("Error fetching %q: %v", u, err)
			continue
		}
		pathComponents := strings.Split(artifactURL.Path, "/")
		fileName := pathComponents[len(pathComponents)-1]

		wg.Add(1)
		go func(urlParam *url.URL, name string) {
			defer wg.Done()
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Transport: tr}
			response, err := client.Get(urlParam.String())
			if err != nil {
				log.Printf("Error fetching %q: %v", urlParam.String(), err)
				return
			}
			defer response.Body.Close()
			r, err := gzip.NewReader(response.Body)
			var b bytes.Buffer
			if _, err := io.Copy(&b, r); err != nil {
				log.Printf("error: %v", err)
			}
			if err := r.Close(); err != nil {
				log.Printf("error: %v", err)
			}
			result[fileName] = b.Bytes()
		}(artifactURL, fileName)
	}
	wg.Wait()
	return result
}

var matchContainerNames = map[string][]string{
	"kube-apiserver":           {"openshift-kube-apiserver_openshift-kube-apiserver"},
	"kube-apiserver-installer": {"openshift-kube-apiserver_installer"},
	"kube-apiserver-operator":  {"openshift-cluster-kube-apiserver-operator_openshift-cluster-kube-apiserver-operator"},

	"kube-controller-manager":           {"openshift-kube-controller-manager_openshift-kube-controller-manager"},
	"kube-controller-manager-installer": {"openshift-kube-controller-manager_installer"},
	"kube-controller-manager-operator":  {"openshift-cluster-kube-controller-manager-operator_openshift-cluster-kube-controller-manager-operator"},

	"openshift-apiserver":          {"openshift-apiserver_apiserver"},
	"openshift-apiserver-operator": {" openshift-cluster-openshift-apiserver-operator_openshift-cluster-openshift-apiserver-operator"},

	"openshift-controller-manager":          {"openshift-controller-manager_controller-manager"},
	"openshift-controller-manager-operator": {" openshift-cluster-openshift-controller-manager-operator_openshift-cluster-openshift-controller-manager-operator"},
}

func LogsCommand(_ Context, jobURL string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	u := parseJobURL(jobURL)
	response, err := client.Get(u)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("get %s failed: HTTP status %d", u, response.StatusCode)
	}

	componentList := []string{}
	for component := range matchContainerNames {
		componentList = append(componentList, component)
	}

	component := ""
	prompt := &survey.Select{
		Message:  "What component you want to see?",
		PageSize: 20,
		Options:  componentList,
	}
	if err := survey.AskOne(prompt, &component, nil); err != nil {
		return err
	}
	if _, exists := matchContainerNames[component]; !exists {
		return fmt.Errorf("invalid component")
	}

	filteredURLs := []string{}
	for _, u := range parseArtifactsURLs(response.Body) {
		for _, m := range matchContainerNames[component] {
			if strings.Contains(u, m) {
				filteredURLs = append(filteredURLs, u)
				break
			}
		}
	}

	for name, body := range fetchArtifacts(filteredURLs) {
		fmt.Printf("\n[%s[:\n", name)
		fmt.Printf("%s\n", string(body))
	}

	return nil
}
