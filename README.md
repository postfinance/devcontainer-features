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
| [docker-out](./features/src/docker-out/README.md) | A feature which installs the Docker client and re-uses the host socket. |
| [go](./features/src/go/README.md) | A feature which installs Go. |
| [zig](./features/src/zig/README.md) | A feature which installs Zig. |

## Usage

### Versions

Most features allow you to define one or more versions of the software that should be installed by the feature.

The preferred way to do this is by defining the full version number, e.g. 1.24.3 for Go. The features by default try to directly download this version and fail if this does not work.

Alternatively, the version often can be set to `latest` or a partial version (e.g. `1.24`). In this case, the feature will resolve the version and find out the appropriate one to install and install that.

NOTE: If you use partial versions, there is also a matching `versionResolve` option that needs to be set to `true` in order for the resolving to be used. This is because some tools do not respect semver and deploy eg. version 2.13.0 as 2.13 instead of 2.13.0, so with this flag, we can distinguish those cases.

NOTE2: The URLs for resolving the versions are usually in separate configuration options. This is so that a feature can be configured to use a cache service (e.g. Artifactory) for downloading the binaries but use the `live` URL for checking versions.

### Global overwrites

Each feature that needs to download something provides options to overwrite the download URL.
This is good if a few projects need a few features.
But if many projects need many features, it can become a nightmare to maintain that.

For this reason, there is a possibility to globally set those overwrites (e.g. for the whole company):

An environment variable `DEV_FEATURE_OVERRIDE_LOCATION` can be set to a location where a text file with the overwrites can be found.
* The variable itself can be defined in the `devcontainer.json` file or already be set in your base images used for dev containers.
* This file can either be on a reachable web or file path. So it can be hosted in a git repository or directly copied into your base images used for dev containers.

The content of the file is simple `key=value` like an env file.
The key names are `DEV_FEATURE_OVERRIDE_<key-to-overwrite>`, but you can also just skip the `DEV_FEATURE_OVERRIDE_` and directly use the desired key name. So for example:
```
GO_DOWNLOAD_URL_BASE=https://mycompany.com/artifactory/dl-google-generic-remote
GO_DOWNLOAD_URL_PATH=/go
```
As the remote can be configured differently (e.g. by including or excluding sub-paths), there are usually two variables: one for the base and one for the path.

Example:

If your remote points to `https://dl.google.com`, you need to set:
* base = "your-remote"
* url = "/go" (the default) or leave out the URL, then it will use the default

If your remote points directly to `https://dl.google.com/go`, you need to set:
* base = "your-remote"
* url = "/" (so it will not add any sub-path after the remote)

See [override-all.env](./override-all.env) for a file with all possible override variables.

### Extend an existing feature

TBD
