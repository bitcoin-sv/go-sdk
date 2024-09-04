# Steps for creating a new release

1. `merge tagged branch to master branch`
2. `git tag v1.0.0`
3. `git push origin v1.0.0`
4. `GITHUB_TOKEN=xxxxxxxx goreleaser release --clean`
