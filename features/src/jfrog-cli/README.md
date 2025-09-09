# JFrog CLI (jfrog-cli)

A package which installs the JFrog CLI.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/jfrog-cli:0.1.0": {
        "version": "latest",
        "downloadUrl": "",
        "versionsUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of JFrog-CLI to install. | string | latest | latest, 2, 2.70, 2.63.1 |
| downloadUrl | The download URL to use for JFrog CLI binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/jfrog-generic-remote |
| versionsUrl | The URL to use to check for available JFrog CLI versions. | string | &lt;empty&gt; |  |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading and resolving the version:
* https://releases.jfrog.io/
