package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

const workers = 8

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repos := make(chan *github.Repository, 2*workers)
	go func() {
		listRepos(ctx, client, repos)
		close(repos)
	}()

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			worker(repos)
		}()
	}

	wg.Wait()
}

func worker(repos chan *github.Repository) {
	for repo := range repos {
		if err := process(repo); err != nil {
			log.Printf("%s: %v\n", repo.GetFullName(), err)
		}
	}
}

func listRepos(ctx context.Context, client *github.Client, repos chan *github.Repository) {
	var opts github.RepositoryListOptions
	for {
		tr, resp, err := client.Repositories.List(ctx, "", &opts)
		if err != nil {
			log.Fatal("Listing repositories:", err)
		}
		for _, repo := range tr {
			repos <- repo
		}
		if resp.NextPage <= opts.Page {
			break
		}
		opts.Page = resp.NextPage
	}
}

func process(repo *github.Repository) error {
	name := repo.GetFullName()
	org := path.Dir(name)
	if err := os.MkdirAll(org, 0o755); err != nil {
		return fmt.Errorf("creating organization dir: %w", err)
	}

	path := name + ".git"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return clone(repo, path)
	}
	return fetch(repo, path)
}

func clone(repo *github.Repository, path string) error {
	log.Println("Clone into", path)
	cmd := exec.Command("git", "clone", "--mirror", repo.GetCloneURL(), path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%s: %s\n", repo.GetFullName(), out)
	}
	return err
}

func fetch(repo *github.Repository, path string) error {
	log.Println("Fetch in", path)
	cmd := exec.Command("git", "remote", "update", "-p")
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%s: %s\n", repo.GetFullName(), out)
	}
	return err
}
