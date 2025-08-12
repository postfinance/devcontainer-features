# devcontainer-features

This repository provides features that can be used in dev containers.

The features provided here have the following benefits:
* AMD64 and ARM compatible (when possible)
* Support for Debian, Ubuntu, Alpine (where possible)
* Corporate ready (download urls can be overwritten with eg. an Artifactory cache but default to the official source)
* Small and fast
* Good extendability

## Features

Below is a list with included features, click on the link for more details.

| Name | Description |
| --- | --- |
| [go](./features/src/go/README.md) | A feature which installs Go. |

## Usage

### Global overwrites

Each feature that needs to download something provides options to overwrite the download url.
This is good if a few projects need a few features.
But if many project need many features, it can get a nightmare to maintain that.

For this reason, there is a possibility to globally set those overwrites (eg. for the whole company):

An environment variable `DEV_FEATURE_OVERRIDE_LOCATION` can be set to a location where a text file with the overwrites can be found.
* The variable itself can be defined in the `devcontainer.json` file or already be set in your base images used for dev containers.
* This file can either be on a reachable web or file path. So it can be hosted in a git repository or directly copied into your base images used for dev containers.

The content if the file is simple `key=value` like an env file.
The key names are `DEV_FEATURE_OVERRIDE_<key-to-overwrite>` but you can also just skip the `DEV_FEATURE_OVERRIDE_` and just directly use the desired key name. So for example:
```
GO_DOWNLOAD_REGISTRY_BASE=https://mycompany.com/artifactory/dl-google-generic-remote
GO_DOWNLOAD_REGISTRY_PATH=/go
```
As the remote can be configured differently (eg. by including or excluding sub-paths), there are usually two variables: one for the base and one for the path.

Example:

Your remote points to `https://dl.google.com`, you need to set:
* base = "your-remote"
* url = "/go" (the default) or leave out the url, then it will use the default

If your remote points directly to `https://dl.google.com/go`, you need to set:
* base = "your-remote"
* url = "/" (so it will not add any subpath after the remote)

### Extend an existing feature

TBD
