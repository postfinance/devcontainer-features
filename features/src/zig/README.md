# Zig (zig)

A feature which installs Zig.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/zig:0.1.0": {
        "version": "latest",
        "versionResolve": false,
        "downloadUrlBase": "",
        "downloadUrlPath": "",
        "versionsUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Zig to install. | string | latest | latest, 0.13.0, 0.12 |
| versionResolve | Whether to resolve the version automatically. | boolean | false | true, false |
| downloadUrlBase | The download URL to use for Zig binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/ziglang-generic-remote |
| downloadUrlPath | The download URL path to use for Zig binaries. | string | &lt;empty&gt; |  |
| versionsUrl | The URL to fetch the available Zig versions from. | string | &lt;empty&gt; |  |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading and resolving:
* https://ziglang.org
