# Docker outside Docker (docker-out)

A feature which installs the Docker client and re-uses the host socket.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/docker-out:0.1.0": {
        "version": "latest",
        "composeVersion": "latest",
        "buildxVersion": "latest"
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of the Docker CLI to install. | string | latest | latest, 28.3.3!, 20.10 |
| composeVersion | The version of the Compose plugin to install. | string | latest | latest, 2.39.1!, 2.29 |
| buildxVersion | The version of the buildx plugin to install. | string | latest | latest, 0.26.1!, 0.10 |

## Customizations

### VS Code Extensions

- `ms-azuretools.vscode-docker`
