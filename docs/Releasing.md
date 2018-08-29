# Releasing a new version

(Notes mostly for myself)

For a new release to get pushed to Github and Docker Hub the following steps
are needed. Set the version you want to release and then run release from the
Makefile like that: `VERSION=0.1.4 make release`.

Then run the following steps. Make sure the Github Token (`GITHUB_TOKEN`) is set for goreleaser.

```
feedbridge|master⚡ ⇒ make release
feedbridge|master⚡ ⇒ make image
feedbridge|master⚡ ⇒ make image-push
```