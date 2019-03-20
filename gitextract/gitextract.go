package gitextract

import (
	"github.com/Masterminds/semver"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"sort"
)

// GetList returns a sorted version list of existing git version tags in given repo
func GetList(repoPath string) (versionList semver.Collection, err error) {
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return versionList, err
	}

	// open repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return versionList, err
	}

	// read tags
	tags, err := repo.Tags()
	if err != nil {
		return versionList, err
	}

	// parse tags
	versionList = semver.Collection{}
	err = tags.ForEach(func(reference *plumbing.Reference) error {
		v, err := semver.NewVersion(reference.Name().Short())
		if err != nil {
			return nil
		}

		versionList = append(versionList, v)

		return nil
	})
	if err != nil {
		return versionList, err
	}

	sort.Sort(versionList)

	return
}
