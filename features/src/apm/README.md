# APM (Agent Package Manager) (apm)

Installs APM (Agent Package Manager) from https://github.com/microsoft/apm/releases.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/apm:1.0.0": {
        "version": "latest",
        "downloadUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of APM to install. | string | latest | latest, 0.26.0, 0.25.0 |
| downloadUrl | The download URL to use for APM binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/github-releases-remote |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading:
* https://github.com

Needs access to the following URL for resolving:
* https://api.github.com
