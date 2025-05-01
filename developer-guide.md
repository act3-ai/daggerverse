# Developer Guide

## Required Tools

Installation Guides:

- [docker](https://docs.docker.com/engine/install/)
- [dagger](https://docs.dagger.io/install)
- [taskfile](https://taskfile.dev/installation/)

## Initializing a New Module

```console
mkdir MODULE_NAME
cd MODULE_NAME
dagger init --sdk=go --name MODULE_NAME
```

## Updating All Module Dependencies

Run `task update-modules`.

## Release Process

TODO