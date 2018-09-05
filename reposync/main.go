package main

import (
	"errors"
	"log"
	"strings"
	"sync"

	git "gopkg.in/src-d/go-git.v4"
)

type LocalRepository struct {
	Name strin
	Path string

	worktree   *git.Worktree
	repository *git.Repository

	err error
}

func (r *LocalRepository) FetchUpstream() {
	if r.hasError() {
		return
	}
	//r.report(r.worktree.Checkout(&git.CheckoutOptions{Branch: "master"}))
	r.report(r.worktree.Pull(&git.PullOptions{RemoteName: "upstream"}))
}

func (r *LocalRepository) report(err error) *LocalRepository {
	if err == nil || r.err != nil {
		return r
	}
	r.err = errors.New(r.Name + ": " + err.Error())
	return r
}

func (r *LocalRepository) hasError() bool {
	return r.err != nil
}

type RepositoryManager struct {
	Repositories []LocalRepository

	err error
}

func (r *RepositoryManager) Add(name, path string) *RepositoryManager {
	if r.hasError() {
		return r
	}
	var err error
	repo := LocalRepository{Name: name, Path: path}

	repo.repository, err = git.PlainOpen(path)
	if err != nil {
		return r.report(err)
	}
	repo.worktree, err = repo.repository.Worktree()
	if err != nil {
		return r.report(err)
	}

	r.Repositories = append(r.Repositories, repo)
	return r
}

func (r *RepositoryManager) Error() error {
	if r.err != nil {
		return r.err
	}
	errs := []string{}
	for _, repo := range r.Repositories {
		if repo.hasError() {
			errs = append(errs, repo.err.Error())
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.New("\n" + strings.Join(errs, "\n"))
}

func (r *RepositoryManager) SynchronizeRepositories() {
	var wg sync.WaitGroup
	for i := range r.Repositories {
		wg.Add(1)
		go func(repository *LocalRepository) {
			defer wg.Done()
			repository.FetchUpstream()
		}(&r.Repositories[i])
	}
	wg.Wait()
}

func (r *RepositoryManager) report(err error) *RepositoryManager {
	r.err = err
	return r
}

func (r *RepositoryManager) hasError() bool {
	return r.err != nil
}

func main() {
	m := &RepositoryManager{}
	m.Add("origin", "/Users/mfojtik/go/src/github.com/openshift/origin").
		Add("k8s.io", "/Users/mfojtik/workspaces/kubernetes/src/k8s.io/kubernetes")

	m.SynchronizeRepositories()

	if m.Error() != nil {
		log.Fatalf("error synchronizing repositories: %v", m.Error())
	}
}
