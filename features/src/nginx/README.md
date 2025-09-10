# Nginx (nginx)

A package which installs Nginx.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/nginx:0.1.0": {
        "version": "latest",
        "stableOnly": false,
        "downloadUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Nginx to install. | string | latest | latest, 1.27.2, 1.27.2-1, 1.26 |
| stableOnly | A flag to indicate if only stable versions should be used. | boolean | false | true, false |
| downloadUrl | The download URL to use for Nginx binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/nginx-generic-remote |

## Notes

The URL in `downloadUrl`, if set, is also used to resolve the version.

### System Compatibility

Debian, Ubuntu

### Accessed Urls

Needs access to the following URL for downloading and resolving the version:
* https://nginx.org
