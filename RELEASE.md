# Steps for creating a new release

1. `git tag v1.0.0`
2. `git push origin v1.0.0`
3. `marge tagged branch to master branch`
4. `goreleaser release`