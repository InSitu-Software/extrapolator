# Extrapolator
Extrapolator provides a go-only helper tool for version / version tag handling, similar to `npm version` or `yarn version`. 

Many thanks to [SemVer](github.com/Masterminds/semver) and [go-git](github.com/src-d/go-git). 

# Usage

    -b, --branch string             branch to be used for commit and tag (default "master")
        --git-tag                   add git tag to current commit
        --major                     increment major version
    -m, --minimal                   minimal output - only print new version string
        --minor                     increment minor version
    -v, --original-version string   use given version as source
        --patch                     increment patch version
        --prefix                    use prefix
        --prefix-tag string         prefix to be added (default "v")
    -r, --repository string         use git repository as source
        --suffix                    use suffix
        --suffix-tag string         suffix to be added (default "beta")

## Example output
    $ extrapolator -r ~/git/exampleRepo --patch --prefix                   
    Current Version: v0.0.47
    New Version: v0.0.48
