# devcontainer-features

This repository provides features that can be used in dev containers.

The features provided here have the following benefits:
* AMD64 and ARM compatible (when possible)
* Support for Debian, Ubuntu, Alpine (where possible)
* Corporate ready (download URLs can be overwritten with e.g. an Artifactory cache but default to the official source)
* Small and fast
* Good extensibility

## Features

Below is a list with included features, click on the link for more details.

| Name | Description |
| --- | --- |
| [browsers](./features/src/browsers/README.md) | A package which installs various browsers. |
| [build-essential](./features/src/build-essential/README.md) | A package which installs build essentials like gcc. |
| [cypress-deps](./features/src/cypress-deps/README.md) | Installs all system dependencies required for running Cypress tests in a dev container. |
| [docker-out](./features/src/docker-out/README.md) | A feature which installs the Docker client and re-uses the host socket. |
| [eclipse-deps](./features/src/eclipse-deps/README.md) | Installs all system dependencies required for running Eclipse IDE in a dev container. |
| [git-lfs](./features/src/git-lfs/README.md) | A feature which installs Git LFS. |
| [go](./features/src/go/README.md) | A feature which installs Go. |
| [goreleaser](./features/src/goreleaser/README.md) | A package which installs GoReleaser. |
| [instant-client](./features/src/instant-client/README.md) | A package which installs the Oracle Instant Client Basic package. |
| [jfrog-cli](./features/src/jfrog-cli/README.md) | A package which installs the JFrog CLI. |
| [locale](./features/src/locale/README.md) | A package which allows setting the locale. |
| [make](./features/src/make/README.md) | A package which installs Make. |
| [mingw](./features/src/mingw/README.md) | A package which installs MinGW. |
| [nginx](./features/src/nginx/README.md) | A package which installs Nginx. |
| [node](./features/src/node/README.md) | A package which installs Node.js. |
| [vault-cli](./features/src/vault-cli/README.md) | A feature which installs the Vault CLI. |
| [zig](./features/src/zig/README.md) | A feature which installs Zig. |

## Usage

### Versions

Most features allow you to define one or more versions of the software that should be installed by the feature.

The logic in the feature to decide which version to install is as follows:
- The passed version number is a full version (e.g. 1.24.3) -> Directly install it and fail if it does not work
- The passed version number is a partial version (e.g. 1.24) -> Lookup the versions and find the highest one which still matches the given one
- The passed version number is `latest` -> Lookup the versions and find the highest one.

The preferred way to do this is always by defining the full version number.

NOTE: The URLs for resolving the versions are usually in separate configuration options. This is so that a feature can be configured to use a cache service (e.g. Artifactory) for downloading the binaries but use the `live` URL for checking versions.

### Global overrides

Each feature that needs to download something provides options to override the download URL.
This is good if a few projects need a few features.
But if many projects need many features, it can become a nightmare to maintain that.

For this reason, there is a possibility to globally set those overrides (e.g. for the whole company):

An environment variable `DEV_FEATURE_OVERRIDE_LOCATION` can be set to a location where a text file with the overrides can be found.
* The variable itself can be defined in the `devcontainer.json` file or already be set in your base images used for dev containers.
* This file can either be on a reachable web or file path. So it can be hosted in a git repository or directly copied into your base images used for dev containers.

The content of the file is simple `key=value` like an env file.
The key names are `DEV_FEATURE_OVERRIDE_<key-to-override>`, but you can also just skip the `DEV_FEATURE_OVERRIDE_` and directly use the desired key name. So for example:
```
GO_DOWNLOAD_URL=https://mycompany.com/artifactory/dl-google-generic-remote/go
```
As the remote can be configured differently (e.g. by including or excluding sub-paths), you might need to check and see which part of the path is expected to be included. Usually as much as possible.

See [override-all.env](./override-all.env) for a file with all possible override variables.

#### Special overrides

There are a few sources which are used in multiple installations. For those sources, there is an override that globaly overrides all installations from this sources. Here is the list of those sources and their keys.

```
DEV_FEATURE_OVERRIDE_GITHUB_DOWNLOAD_URL=...
```

### Extend an existing feature

TBD
