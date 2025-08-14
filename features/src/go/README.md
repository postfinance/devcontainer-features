# Go (go)

A feature which installs Go.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/go:0.1.0": {
        "version": "latest",
        "versionResolve": false,
        "downloadRegistryBase": "",
        "downloadRegistryPath": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Go to install. | string | latest | latest, 1.24, 1.21.8 |
| versionResolve | Whether to resolve the version automatically. | boolean | false | true, false |
| downloadRegistryBase | The download registry to use for Go binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/dl-google-generic-remote |
| downloadRegistryPath | The download registry path to use for Go binaries. | string | &lt;empty&gt; |  |

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
