package gitextract

import (
	"github.com/Masterminds/semver"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	plumberSSH "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

// GetMaxVersionByLocal returns the latest / max version of a given / cloned git repo
func GetMaxVersionByLocal(repoPath string) (maxVersion *semver.Version, err error) {
	list, err := GetListByLocal(repoPath)
	if err != nil {
		return
	}

	return getMaxVersion(list)
}

// GetMaxVersionByRemote returns the latest / max version of a given git repo by url
func GetMaxVersionByRemote(repoURL, key, user, password string) (maxVersion *semver.Version, err error) {
	list, err := GetListByRemote(repoURL, key, user, password)
	if err != nil {
		return
	}

	return getMaxVersion(list)
}

func getMaxVersion(list semver.Collection) (maxVersion *semver.Version, err error) {
	if len(list) == 0 {
		maxVersion, err = semver.NewVersion("0.0.0")
		return
	}

	sort.Sort(list)
	maxVersion = list[len(list)-1]

	return maxVersion, nil
}

func getVersionList(repo *git.Repository) (versionList semver.Collection, err error) {
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

	return
}

// GetListByLocal returns a list of existing git version tags in given local repo
func GetListByLocal(repoPath string) (versionList semver.Collection, err error) {
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return versionList, err
	}

	// open repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return versionList, err
	}

	versionList, err = getVersionList(repo)

	return
}

// GetListByRemote returns a list of existing git version tags in given remote repo
func GetListByRemote(remote, key, user, password string) (versionList semver.Collection, err error) {

	cloneOptions := git.CloneOptions{
		URL: remote,
	}

	switch {
	case key != "":
		privateKey, err := ioutil.ReadFile(key)
		if err != nil {
			return versionList, err
		}

		signer, err := ssh.ParsePrivateKey(privateKey)
		if err != nil {
			return versionList, err
		}

		cloneOptions.Auth = &plumberSSH.PublicKeys{User: user, Signer: signer}
	case password != "":
		log.Fatal("not implemented")
	case user != "git":
		log.Fatal("not implemented")
	}

	repo, err := git.Clone(memory.NewStorage(), nil, &cloneOptions)
	if err != nil {
		return
	}

	versionList, err = getVersionList(repo)
	return
}
