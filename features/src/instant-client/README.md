# Instant client (instant-client)

Installs the Oracle Instant Client Basic package.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/instant-client:0.1.0": {
        "version": "latest",
        "downloadUrl": "",
        "versionsUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Instant Client to install. | string | latest | latest, 23, 21, 23.8.0.25.04 |
| downloadUrl | The download URL to use for Instant Client binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/oracle-generic-remote |
| versionsUrl | The URL to use to check for available Instant Client versions. | string | &lt;empty&gt; |  |

## Notes

Restrictions:
* Versions **below 19** are not supported! 
* Version 21 is not supported for ARM!

### System Compatibility

Debian, Ubuntu