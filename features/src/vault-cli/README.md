# Vault CLI (vault-cli)

A feature which installs the Vault CLI.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/vault-cli:0.2.0": {
        "version": "latest",
        "downloadUrl": "",
        "versionsUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Vault CLI to install. | string | latest | latest, 1.18.2, 1 |
| downloadUrl | The download URL to use for Vault CLI binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/hashicorp-releases-generic-remote/vault |
| versionsUrl | The URL to fetch the available Vault CLI versions from. | string | &lt;empty&gt; |  |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading and resolving:
* https://releases.hashicorp.com
