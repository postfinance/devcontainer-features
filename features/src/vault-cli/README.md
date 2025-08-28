# Vault CLI (vault-cli)

A feature which installs the Vault CLI.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/vault-cli:0.1.0": {
        "version": "latest",
        "versionResolve": false,
        "downloadUrlBase": "",
        "downloadUrlPath": "",
        "versionsUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Vault CLI to install. | string | latest | latest, 1.18.2, 1 |
| versionResolve | Whether to resolve the version automatically. | boolean | false | true, false |
| downloadUrlBase | The download URL to use for Vault CLI binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/hashicorp-releases-generic-remote |
| downloadUrlPath | The download URL path to use for Vault CLI binaries. | string | &lt;empty&gt; |  |
| versionsUrl | The URL to fetch the available Vault CLI versions from. | string | &lt;empty&gt; |  |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading and resolving:
* https://releases.hashicorp.com
