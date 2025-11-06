# Go (go)

Installs Go.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/go:0.2.0": {
        "version": "latest",
        "downloadUrl": "",
        "latestUrl": "",
        "versionsUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Go to install. | string | latest | latest, 1.24, 1.21.8 |
| downloadUrl | The download URL to use for Go binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/dl-google-generic-remote/go |
| latestUrl | The URL to fetch the latest Go version from. | string | &lt;empty&gt; |  |
| versionsUrl | The URL to fetch the available Go versions from. | string | &lt;empty&gt; |  |

## Customizations

### VS Code Extensions

- `golang.Go`

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading:
* https://dl.google.com

Needs access to the following URL for resolving:
* https://go.dev
