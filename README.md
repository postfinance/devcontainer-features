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
| [browsers](./features/src/browsers/README.md) | Installs various browsers and their dependencies. |
| [build-essential](./features/src/build-essential/README.md) | Installs build essentials like gcc. |
| [cypress-deps](./features/src/cypress-deps/README.md) | Installs all dependencies required to run Cypress. |
| [docker-out](./features/src/docker-out/README.md) | Installs a Docker client which re-uses the host Docker socket. |
| [eclipse-deps](./features/src/eclipse-deps/README.md) | Installs all dependencies required to run the Eclipse IDE. |
| [git-lfs](./features/src/git-lfs/README.md) | Installs Git LFS. |
| [gitlab-cli](./features/src/gitlab-cli/README.md) | Installs the GitLab CLI. |
| [go](./features/src/go/README.md) | Installs Go. |
| [gonovate](./features/src/gonovate/README.md) | Installs Gonovate. |
| [goreleaser](./features/src/goreleaser/README.md) | Installs GoReleaser. |
| [instant-client](./features/src/instant-client/README.md) | Installs the Oracle Instant Client Basic package. |
| [jfrog-cli](./features/src/jfrog-cli/README.md) | Installs the JFrog CLI. |
| [locale](./features/src/locale/README.md) | Allows setting the locale. |
| [make](./features/src/make/README.md) | Installs Make. |
| [mingw](./features/src/mingw/README.md) | Installs MinGW. |
| [nginx](./features/src/nginx/README.md) | Installs Nginx. |
| [node](./features/src/node/README.md) | Installs Node.js. |
| [playwright-deps](./features/src/playwright-deps/README.md) | Installs all dependencies required to run Playwright. |
| [python](./features/src/python/README.md) | Installs Python. |
| [sonar-scanner-cli](./features/src/sonar-scanner-cli/README.md) | Installs the SonarScanner CLI. |
| [system-packages](./features/src/system-packages/README.md) | Install arbitrary system packages using the system package manager. |
| [timezone](./features/src/timezone/README.md) | Allows setting the timezone. |
| [vault-cli](./features/src/vault-cli/README.md) | Installs the Vault CLI. |
| [zig](./features/src/zig/README.md) | Installs Zig. |

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

The precedence for the overrides is:

1. Value set via feature parameter
2. Value set via environment variable
3. Values from `DEV_FEATURE_OVERRIDE_LOCATION`

#### Special overrides

There are a few sources which are used in multiple installations. For those sources, there is an override that globaly overrides all installations from this sources. Here is the list of those sources and their keys.

```
DEV_FEATURE_OVERRIDE_GITHUB_DOWNLOAD_URL=...
```
#### Unset an Override via Parameter

If an override is set, setting the corresponding parameter to `""` will not unset the override. To achieve this, set the parameter to `none`.

**Example:**  
This environment variable is set: `DOCKER_OUT_CONFIG_PATH=https://example.com/config.json`

Then set this in your feature to explicitly unset it:

```json
{
    "ghcr.io/postfinance/devcontainer-features/docker-out:0.3.0": {
        "version": "28.3.3",
        "configPath": "none"
    }
}
```

### Extend an existing feature

TBD
