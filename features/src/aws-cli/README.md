# AWS CLI (aws-cli)

Installs the AWS CLI.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/aws-cli:1.0.0": {
        "version": "latest",
        "downloadUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of AWS CLI to install. | string | latest | latest, 2.22.35 |
| downloadUrl | The download URL to use for AWS CLI binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/awscli-remote |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine

### Accessed Urls

Needs access to the following URL for downloading and resolving:
* https://awscli.amazonaws.com
* https://github.com
* https://api.github.com
