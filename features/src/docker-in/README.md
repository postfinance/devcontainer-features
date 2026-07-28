# Docker inside Docker (docker-in)

Installs and runs a full Docker daemon.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/docker-in:1.0.0": {
        "version": "latest",
        "composeVersion": "latest",
        "buildxVersion": "latest",
        "configPath": "",
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
| version | The version of the Docker Engine and CLI to install. | string | latest | latest, 28.3.3, 20.10 |
| composeVersion | The version of the Compose plugin to install. | string | latest | latest, none, 2.39.1, 2.29 |
| buildxVersion | The version of the buildx plugin to install. | string | latest | latest, none, 0.26.1, 0.10 |
| configPath | Path or URL to a custom Docker client config.json file to copy into the container. | string | &lt;empty&gt; | /home/user/.docker/config.json, https://raw.githubusercontent.com/devcontainers/features/main/src/docker-in/config.json, none |
| downloadUrl | The download URL to use for Docker binaries (engine, CLI, containerd). | string | &lt;empty&gt; |  |
| versionsUrl | The URL to use for checking available versions. | string | &lt;empty&gt; |  |
| composeDownloadUrl | The download URL to use for Docker Compose binaries. | string | &lt;empty&gt; |  |
| buildxDownloadUrl | The download URL to use for Docker Buildx binaries. | string | &lt;empty&gt; |  |

## Customizations

### VS Code Extensions

- `ms-azuretools.vscode-containers`
