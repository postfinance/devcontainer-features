# Waza (waza)

Installs Waza, the AI-powered development agent from Microsoft.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/waza:1.0.0": {
        "version": "latest",
        "downloadUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Waza to install. | string | latest | latest, 0.38.5 |
| downloadUrl | The download URL to use for Waza binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/github-releases-remote |
