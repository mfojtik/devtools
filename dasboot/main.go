package main

import (
	"log"
	"net/http"
	"time"
)

var componentURLs = map[string]string{
	"kube-apiserver-operator-e2e-aws":               "https://openshift-gce-devel.appspot.com/builds/origin-ci-test/pr-logs/directory/pull-ci-openshift-cluster-kube-apiserver-operator-master-e2e-aws",
	"kube-controller-manager-operator-e2e-aws":      "https://openshift-gce-devel.appspot.com/builds/origin-ci-test/pr-logs/directory/pull-ci-openshift-cluster-kube-controller-manager-operator-master-e2e-aws",
	"openshift-apiserver-operator-e2e-aws":          "https://openshift-gce-devel.appspot.com/builds/origin-ci-test/pr-logs/directory/pull-ci-openshift-cluster-openshift-apiserver-operator-master-e2e-aws",
	"openshift-controller-manager-operator-e2e-aws": "https://openshift-gce-devel.appspot.com/builds/origin-ci-test/pr-logs/directory/pull-ci-openshift-cluster-openshift-controller-manager-operator-master-e2e-aws",
}

type Status string

var (
	StatusSuccess Status = "Success"
	StatusFailed  Status = "Failed"
	StatusPending Status = "Pending"
)

type Build struct {
	Status    Status    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type Component struct {
	Name   string           `json:"name"`
	Builds map[string]Build `json:"builds"`
}

var serveMux = http.NewServeMux()

func main() {
	stopCh := make(chan struct{})
	defer close(stopCh)

	go newPeriodicScraper(5*time.Minute, stopCh)

	serveMux.Handle("/", http.FileServer(http.Dir("./html/")))
	http.HandleFunc("/", statusHandler)

	log.Printf("Serving status at http://0.0.0.0:8080 ...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("error listening: %v", err)
	}
}
