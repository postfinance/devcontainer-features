# GitHub Copilot CLI (github-copilot-cli)

Installs GitHub Copilot CLI (copilot), the AI-powered coding assistant for the terminal.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/github-copilot-cli:1.0.0": {
        "version": "latest",
        "downloadUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of GitHub Copilot CLI to install. | string | latest | latest, 1.0.39 |
| downloadUrl | The download URL to use for GitHub Copilot CLI binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/github-releases-remote |

## Notes

### System Compatibility

Debian, Ubuntu

### Accessed Urls

Needs access to the following URL for downloading:
* https://github.com

Needs access to the following URL for resolving:
* https://api.github.com
