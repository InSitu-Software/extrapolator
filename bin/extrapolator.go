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
	pflag.StringP("branch", "b", "master", "branch to be used for commit and tag")
	pflag.StringP("original-version", "v", "", "use given version as source")

	// define output
	pflag.Bool("major", false, "increment major version")
	pflag.Bool("minor", false, "increment minor version")
	pflag.Bool("patch", false, "increment patch version")

	pflag.Bool("suffix", false, "use suffix")
	pflag.String("suffix-tag", "beta", "suffix to be added")

	pflag.Bool("prefix", false, "use prefix")
	pflag.String("prefix-tag", "v", "prefix to be added")

	// define actions
	pflag.Bool("git-tag", false, "add git tag to current commit")

	// output options
	pflag.BoolP("minimal", "m", false, "minimal output - only print new version string")

	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal(err)
	}

	if viper.GetString("repository") != "" && viper.GetString("original-version") != "" {
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
	switch {
	case viper.GetString("original-version") != "":
		v, err := semver.NewVersion(viper.GetString("original-version"))
		if err != nil {
			log.Fatal(err)
		}
		maxVersion = v
	case viper.GetString("repository") != "":
		list, err := gitextract.GetList(viper.GetString("repository"))
		if err != nil {
			log.Fatal(err)
		}
		if len(list) == 0 {
			if maxVersion, err = semver.NewVersion("0.0.0"); err != nil {
				log.Fatal(err)
			}
			break
		}
		maxVersion = list[len(list)-1]

	default:
		log.Fatal("missing datasource")
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
	if viper.GetBool("git-tag") {
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
