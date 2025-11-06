# Git LFS (git-lfs)

Installs Git LFS.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/git-lfs:0.2.0": {
        "version": "latest",
        "downloadUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Git LFS to install. | string | latest | latest, 3.7.0, 3.6 |
| downloadUrl | The download URL to use for Git LFS binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/github-generic-remote/git-lfs/git-lfs/releases/download, https://mycompany.com/artifactory/git-lfs-download-generic-remote |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading and resolving:
* https://github.com
