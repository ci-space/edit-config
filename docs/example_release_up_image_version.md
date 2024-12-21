# Release workflow with up image version in GitHub action.yaml

<p align="center">
  <img src="./docs/usage_release_1.png" />
  <img src="./docs/usage_release_2.png" />
</p>

**./github/workflows/release.yaml:**

```yaml
name: release

permissions: write-all

on:
  workflow_dispatch:
    inputs:
      version:
        description: version
        required: true
        type: choice
        options:
          - major
          - minor
          - patch

jobs:
  release:
    name: release
    runs-on: ubuntu-latest
    steps:
      - name: Check out code
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4 # action page: <https://github.com/actions/setup-go>
        with:
          go-version: stable

      - name: Generate next version
        id: version
        uses: ci-space/edit-config@master
        with:
          file: action.yaml
          action: up-image-version
          pointer: runs.image
          value: ${{ github.event.inputs.version }}

      - name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Build image
        uses: docker/build-push-action@v4 # Action page: <https://github.com/docker/build-push-action>
        with:
          context: .
          file: Dockerfile
          push: true
          platforms: linux/amd64
          tags: |
            ghcr.io/ci-space/edit-config:${{ steps.version.outputs.new-version }}

      - name: Commit changes
        run: |
          git config user.name github-actions[bot]
          git config user.email github-actions[bot]@users.noreply.github.com
          git add action.yaml
          git commit -m "chore: update image version ${{ steps.version.outputs.new-version }} in action.yaml"
          git push

      - name: Create Tag
        uses: negz/create-tag@v1
        with:
          version: ${{ steps.version.outputs.new-version }}
          message: ''
          token: ${{ secrets.GITHUB_TOKEN }}
```
