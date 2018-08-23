# Releasing a new version

For a new release to get pushed to Github and Docker Hub the following steps
are needed. Set version in `Makefile` and then run the following steps. Make sure
the Github Token (`GITHUB_TOKEN`) is set for goreleaser.

IMPORTANT: Set version in the Makefile before running this.

```
feedbridge|master⚡ ⇒ make image
feedbridge|master⚡ ⇒ make image-push
feedbridge|master⚡ ⇒ make release
```