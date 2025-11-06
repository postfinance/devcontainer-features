# Docker outside Docker (docker-out)

Installs a Docker client which re-uses the host Docker socket.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/docker-out:0.2.0": {
        "version": "latest",
        "composeVersion": "latest",
        "buildxVersion": "latest",
        "downloadUrl": "",
        "versionsUrl": "",
        "composeDownloadUrl": "",
        "buildxDownloadUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of the Docker CLI to install. | string | latest | latest, 28.3.3, 20.10 |
| composeVersion | The version of the Compose plugin to install. | string | latest | latest, none, 2.39.1, 2.29 |
| buildxVersion | The version of the buildx plugin to install. | string | latest | latest, none, 0.26.1, 0.10 |
| downloadUrl | The download URL to use for Docker binaries. | string | &lt;empty&gt; |  |
| versionsUrl | The URL to use for checking available versions. | string | &lt;empty&gt; |  |
| composeDownloadUrl | The download URL to use for Docker Compose binaries. | string | &lt;empty&gt; |  |
| buildxDownloadUrl | The download URL to use for Docker Buildx binaries. | string | &lt;empty&gt; |  |

## Customizations

### VS Code Extensions

- `ms-azuretools.vscode-docker`

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading and resolving:
* https://download.docker.com
* https://github.com
