package main

import (
	"fmt"
	"github.com/InSitu-Software/extrapolator/gitextract"
	"github.com/Masterminds/semver"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
)

func init() {
	// define data source
	pflag.StringP("repository", "r", "", "use git repository as source")
	pflag.String("repository-remote", "", "use a not cloned repository")
	pflag.StringP("branch", "b", "master", "branch to be used for commit and tag")
	pflag.StringP("original-version", "v", "", "use given version as source")

	// define auth for repository-remote if necessary
	pflag.String("remote-ssh", "", "ssh key to be used for auth")
	pflag.String("remote-user", "git", "user to be used for auth")
	pflag.String("remote-password", "", "password to be used for auth")

	// define output
	pflag.Bool("major", false, "increment major version")
	pflag.Bool("minor", false, "increment minor version")
	pflag.Bool("patch", false, "increment patch version")

	pflag.Bool("suffix", false, "use suffix")
	pflag.String("suffix-tag", "beta", "suffix to be added")

	pflag.Bool("prefix", false, "use prefix")
	pflag.String("prefix-tag", "v", "prefix to be added")

	// define actions
	pflag.Bool("git-tag", false, "add git tag to current commit (works only on local repositories")

	// output options
	pflag.BoolP("minimal", "m", false, "minimal output - only print new version string")

	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal(err)
	}

	if viper.GetString("repository") != "" &&
		viper.GetString("original-version") != "" &&
		viper.GetString("repository-remote") == "" {
		log.Fatal("please use original-version OR repository as source")
	}

	// check for only one set parameter to be changed
	var s int
	for _, v := range []bool{viper.GetBool("major"), viper.GetBool("minor"), viper.GetBool("patch")} {
		if v {
			s++
		}
	}
	if s > 1 {
		log.Fatal("please specify exactly one of major, minor, patch zu be incremented")
	}
}

func main() {
	// get latest version
	var maxVersion *semver.Version
	var err error
	switch {
	case viper.GetString("original-version") != "":
		maxVersion, err = semver.NewVersion(viper.GetString("original-version"))
	case viper.GetString("repository") != "":
		maxVersion, err = gitextract.GetMaxVersionByLocal(viper.GetString("repository"))
	case viper.GetString("repository-remote") != "":
		maxVersion, err = gitextract.GetMaxVersionByRemote(
			viper.GetString("repository-remote"),
			viper.GetString("remote-ssh"),
			viper.GetString("remote-user"),
			viper.GetString("remote-password"),
		)
	default:
		log.Fatal("missing datasource")
	}
	if err != nil {
		log.Fatal(err)
	}

	// generate new version
	var newVersion semver.Version
	switch {
	case viper.GetBool("major"):
		newVersion = maxVersion.IncMajor()
	case viper.GetBool("minor"):
		newVersion = maxVersion.IncMinor()
	case viper.Get("patch"):
		newVersion = maxVersion.IncPatch()
	default:
		newVersion = *maxVersion
	}

	// version string preparation
	var newVersionString string
	switch {
	case viper.GetBool("suffix"):
		suffixVersion, err := newVersion.SetPrerelease(viper.GetString("suffix-tag"))
		if err != nil {
			log.Fatal(err)
		}
		newVersionString = suffixVersion.String()
	default:
		newVersionString = newVersion.String()
	}

	if viper.GetBool("prefix") {
		newVersionString = fmt.Sprintf("%s%s", viper.GetString("prefix-tag"), newVersionString)
	}

	// create git tag
	if viper.GetBool("git-tag") && viper.GetString("repository") != "" {
		err := gitextract.Tag(viper.GetString("repository"), newVersionString, viper.GetString("branch"))
		if err != nil {
			log.Fatal(err)
		}
	}

	// output
	if viper.GetBool("minimal") {
		fmt.Println(newVersionString)
		os.Exit(0)
	}

	fmt.Printf("Current Version: %s\nNew Version: %s\n", maxVersion.Original(), newVersionString)
}
