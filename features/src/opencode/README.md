# opencode (opencode)

Installs opencode, the open source AI coding agent.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/opencode:1.0.0": {
        "version": "latest",
        "downloadUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of opencode to install. | string | latest | latest, 1.14.28 |
| downloadUrl | The download URL to use for opencode binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/github-releases-remote |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading:
* https://github.com

Needs access to the following URL for resolving:
* https://api.github.com
