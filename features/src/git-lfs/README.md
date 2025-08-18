# Git LFS (git-lfs)

A feature which installs Git LFS.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/git-lfs:0.1.0": {
        "version": "latest",
        "versionResolve": false,
        "downloadUrlBase": "",
        "downloadUrlPath": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Git LFS to install. | string | latest | latest, 3.7.0, 3.6 |
| versionResolve | Whether to resolve the version automatically. | boolean | false | true, false |
| downloadUrlBase | The download URL to use for Git LFS binaries. | string | &lt;empty&gt; |  |
| downloadUrlPath | The download URL path to use for Git LFS binaries. | string | &lt;empty&gt; |  |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading and resolving:
* https://github.com
