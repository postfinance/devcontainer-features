# Docker outside Docker (docker-out)

A feature which installs the Docker client and re-uses the host socket.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/docker-out:0.1.0": {
        "version": "latest",
        "versionResolve": false,
        "composeVersion": "latest",
        "composeVersionResolve": false,
        "buildxVersion": "latest",
        "buildxVersionResolve": false,
        "downloadUrlBase": "",
        "downloadUrlPath": "",
        "versionsUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of the Docker CLI to install. | string | latest | latest, 28.3.3, 20.10 |
| versionResolve | Whether to resolve the version automatically. | boolean | false | true, false |
| composeVersion | The version of the Compose plugin to install. | string | latest | latest, 2.39.1, 2.29 |
| composeVersionResolve | Whether to resolve the version automatically. | boolean | false | true, false |
| buildxVersion | The version of the buildx plugin to install. | string | latest | latest, 0.26.1, 0.10 |
| buildxVersionResolve | Whether to resolve the version automatically. | boolean | false | true, false |
| downloadUrlBase | The download URL to use for Docker binaries. | string | &lt;empty&gt; |  |
| downloadUrlPath | The download URL path to use for Docker binaries. | string | &lt;empty&gt; |  |
| versionsUrl | The URL to use for checking available versions. | string | &lt;empty&gt; |  |

## Customizations

### VS Code Extensions

- `ms-azuretools.vscode-docker`

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading and resolving:
* https://download.docker.com
