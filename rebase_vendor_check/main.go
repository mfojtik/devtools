package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/mfojtik/devtools/rebase_vendor_check/pkg/godep"
)

type dependency struct {
	path     string
	revision string
}

var notFoundErr = errors.New("dependency not found")

func goGetDep(d dependency) error {
	cmd := exec.Command("go", "get", "-u", d.path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git: %q (%v)", out, err)
	}
	return nil
}

func gitCommitExists(d dependency, commitId string) bool {
	cmd := exec.Command("git", "show", commitId)
	cmd.Dir = filepath.Join(os.Getenv("GOPATH"), "src", d.path)
	if _, err := os.Stat(cmd.Dir); os.IsNotExist(err) {
		return false
	}
	_, err := cmd.CombinedOutput()
	return err == nil
}

func gitFetchDep(d dependency) error {
	cmd := exec.Command("git", "fetch")
	cmd.Dir = filepath.Join(os.Getenv("GOPATH"), "src", d.path)
	if _, err := os.Stat(cmd.Dir); os.IsNotExist(err) {
		return notFoundErr
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git: %q (%v)", out, err)
	}
	return nil
}

func loadDependencies(source godep.Godeps) map[string]dependency {
	compacted := map[string]dependency{}
	for _, d := range source.Deps {
		compactedImport := d.ImportPath
		parts := strings.Split(d.ImportPath, "/")
		if len(parts) > 3 {
			compactedImport = strings.Join(parts[0:3], "/")
		}
		compacted[compactedImport] = dependency{path: compactedImport, revision: d.Rev}
	}
	return compacted
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("usage: %s kubernetes/Godep.json origin/Godep.json # must be run in origin directory", os.Args[0])
	}
	kubeGodeps, err := godep.LoadGodepsFile(os.Args[1])
	if err != nil {
		log.Fatalf("error loading %s: %v", os.Args[1], err)
	}
	originGodeps, err := godep.LoadGodepsFile(os.Args[2])
	if err != nil {
		log.Fatalf("error loading %s: %v", os.Args[2], err)
	}

	compactedKubeDeps := loadDependencies(kubeGodeps)
	compactedOriginDeps := loadDependencies(originGodeps)
	workItems := []dependency{}

	for _, k := range compactedKubeDeps {
		depLog := log.WithField("k", k.revision)
		o, ok := compactedOriginDeps[k.path]
		if !ok {
			depLog.Infof("new: %s", k.path)
			continue
		}
		depLog = depLog.WithField("o", o.revision)
		if o.revision != k.revision {
			depLog.Infof("change: %s", k.path)
			workItems = append(workItems, o)
		}
	}

	fmt.Println()
	log.Infof("Processing changes ...")

	for _, o := range workItems {
		logger := log.WithField("d", o.path)

		err := gitFetchDep(o)
		if err != nil && err != notFoundErr {
			log.Warnf("error fetching dependency: %v", err)
		}

		if err == notFoundErr {
			if err := goGetDep(o); err != nil {
				log.Warnf("unable to go get: %v", err)
			}
		}

		kubeLevel := compactedKubeDeps[o.path].revision
		if !gitCommitExists(o, kubeLevel) {
			logger.Warnf("kube level of %q does not exists", kubeLevel)
			continue
		}
		if !gitCommitExists(o, o.revision) {
			logger.Warnf("origin level of %q does not exists (missing picks?)", o.revision)
			continue
		}
		log.Infof("%s: safe to bump to %q", o.path, kubeLevel)
	}
}
