# JFrog CLI (jfrog-cli)

A package which installs the JFrog CLI.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/jfrog-cli:0.1.0": {
        "version": "latest"
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of JFrog-CLI to install. | string | latest | latest, 2, 2.70, 2.63.1 |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading and resolving the version:
* https://releases.jfrog.io/
