# edit-config

edit-config - is console app and GitHub action for editing configuration files

## Workflow examples based on edit-config
- [Manual release workflow with up image version](./docs/example_release_up_image_version.md)

## Usage

Available actions:
- Up image version
- Append: add value to array / concat string / add number

### Up image version

```yaml
jobs:
  up-image-version:
    steps:
      - name: Up image version
        id: version
        uses: ci-space/edit-config@master
        with:
          file: action.yaml
          action: up-image-version
          pointer: runs.image
          value: ${{ github.event.inputs.version }}

      - name: Print image version
        run: echo ${{ steps.version.outputs.new-version }}
```

### Append value to array

```yaml
jobs:
  append-value-to-array:
    steps:
      - name: Append value to array
        uses: ci-space/edit-config@master
        with:
          file: users.yaml
          action: append
          pointer: users[0].phones
          value: '123'
```
