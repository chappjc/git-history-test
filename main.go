package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func checkIfError(args ...interface{}) {
	if len(args) < 2 || args[1] == nil {
		return
	}
	log.Fatal(args...)
	os.Exit(1)
}

func main() {
	log.Println("git clone https://github.com/dmigwi/golang-modules.git")
	r, err := git.PlainClone("data", false, &git.CloneOptions{
		URL: "https://github.com/dmigwi/golang-modules.git",
	})
	if err == git.ErrRepositoryAlreadyExists {
		r, err = git.PlainOpen("data")
		checkIfError(err)

		var w *git.Worktree
		w, err = r.Worktree()
		checkIfError(err)

		// Pull the latest changes from the origin remote and merge into the current branch
		log.Println("git pull origin")
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	}
	// r, err := git.PlainOpen("../golang-modules")
	checkIfError(err)

	branchHash, err := r.ResolveRevision("master")
	checkIfError(err)

	commits, err := r.Log(&git.LogOptions{From: *branchHash})
	checkIfError(err)

	// Starting at HEAD, commit messages and diffs.
	var futureCommit, commit *object.Commit
	for {
		commit, err = commits.Next()
		if err != nil {
			if err == io.EOF {
				// no more commits
				break
			} else {
				log.Fatal(err)
				os.Exit(1)
			}
		}

		if futureCommit == nil {
			futureCommit = commit
			continue
		}

		diff, err := commit.Patch(futureCommit)
		checkIfError(err)

		fmt.Println(futureCommit.String())
		fmt.Println(diff.String())

		// Continue to the previous commit.
		futureCommit = commit
	}

	os.Exit(0)
}
