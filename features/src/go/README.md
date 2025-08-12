# Go (go)

A feature which installs Go.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/go:1.0.0": {
        "version": "latest",
        "isExactVersion": false,
        "downloadRegistryBase": "",
        "downloadRegistryPath": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Go to install. | string | latest | latest, 1.24, 1.21.8 |
| isExactVersion | Whether to install the exact version specified. | boolean | false | true, false |
| downloadRegistryBase | The download registry to use for Go binaries. | string | <empty> | https://mycompany.com/artifactory/dl-google-generic-remote |
| downloadRegistryPath | The download registry path to use for Go binaries. | string | <empty> |  |

## Customizations

### VS Code Extensions

- `golang.Go`

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs to access the following urls to install non-exact versions:
* https://go.dev/dl/?mode=json&include=all
* https://go.dev/VERSION?m=text
