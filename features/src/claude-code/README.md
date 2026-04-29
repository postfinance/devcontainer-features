# claude-code (claude-code)

Installs Claude Code, Anthropic's AI coding assistant CLI.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/claude-code:1.0.0": {
        "version": "latest",
        "downloadUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Claude Code to install. | string | latest | latest, 2.1.123 |
| downloadUrl | The download URL to use for Claude Code binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/github-releases-remote |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading:
* https://github.com

Needs access to the following URL for resolving:
* https://api.github.com
