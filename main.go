package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	var repos []*github.Repository
	var opts github.RepositoryListOptions
	for {
		tr, resp, err := client.Repositories.List(ctx, "", &opts)
		if err != nil {
			log.Fatal("Listing repositories:", err)
		}
		repos = append(repos, tr...)
		if resp.NextPage <= opts.Page {
			break
		}
		opts.Page = resp.NextPage
	}

	for _, repo := range repos {
		if err := process(repo); err != nil {
			log.Printf("%s: %v\n", repo.GetFullName(), err)
		}
	}
}

func process(repo *github.Repository) error {
	name := repo.GetFullName()
	org := path.Dir(name)
	if err := os.MkdirAll(org, 0755); err != nil {
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
	_, err := cmd.CombinedOutput()
	return err
}

func fetch(repo *github.Repository, path string) error {
	log.Println("Fetch in", path)
	cmd := exec.Command("git", "remote", "update", "-p")
	cmd.Dir = path
	_, err := cmd.CombinedOutput()
	return err
}
