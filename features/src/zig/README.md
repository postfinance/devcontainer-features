# Zig (zig)

A feature which installs Zig.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/zig:1.0.0": {
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
| version | The version of Zig to install. | string | latest | latest, 0.13.0, 0.12 |
| isExactVersion | Whether to install the exact version specified. | boolean | false | true, false |
| downloadRegistryBase | The download registry to use for Zig binaries. | string | <empty> | https://mycompany.com/artifactory/ziglang-generic-remote |
| downloadRegistryPath | The download registry path to use for Zig binaries. | string | <empty> |  |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs to access the following urls to install non-exact versions:
* https://ziglang.org/download/index.json
