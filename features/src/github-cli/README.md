# GitHub CLI (github-cli)

Installs the GitHub CLI.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/github-cli:1.0.0": {
        "version": "latest",
        "downloadUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of GitHub CLI to install. | string | latest | latest, 2.67.0 |
| downloadUrl | The download URL to use for GitHub CLI binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/github-releases-remote |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading and resolving:
* https://github.com
* https://api.github.com
