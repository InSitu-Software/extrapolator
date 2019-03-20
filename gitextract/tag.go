package gitextract

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
)

func Tag(repoPath string, tag string, branch string) (err error) {
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return err
	}

	// open repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return
	}

	// switch to branch
	worktree, err := repo.Worktree()
	if err != nil {
		return
	}

	var checkoutOptions git.CheckoutOptions
	if branch != "" && branch != "master" {
		branchReference, err := repo.Branch("branch")
		if err != nil {
			return err
		}

		checkoutOptions.Branch = plumbing.ReferenceName(branchReference.Name)

		if err := checkoutOptions.Validate(); err != nil {
			return err
		}
	}

	if err := worktree.Checkout(&checkoutOptions); err != nil {
		return err
	}

	headReference, err := repo.Head()
	if err != nil {
		return err
	}

	if _, err := repo.CreateTag(tag, headReference.Hash(), nil); err != nil {
		return err
	}

	return
}
