# GitLab CLI (gitlab-cli)

A package which installs the GitLab CLI.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/gitlab-cli:0.1.0": {
        "version": "latest",
        "downloadUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of GitLab CLI to install. | string | latest | latest, 1.67.0 |
| downloadUrl | The download URL to use for the binaries. | string | &lt;empty&gt; |  |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading and resolving:
* https://gitlab.com
