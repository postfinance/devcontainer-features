# GoReleaser (goreleaser)

A package which installs GoReleaser.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/goreleaser:0.1.0": {
        "version": "latest",
        "downloadUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of GoReleaser to install. | string | latest | latest, 2.4.8, 2.3 |
| downloadUrl | The download URL to use for GoReleaser binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/github-releases-remote |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading:
* https://github.com

Needs access to the following URL for resolving:
* https://api.github.com
