# Terraform (terraform)

Installs the Terraform CLI.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/terraform:0.1.0": {
        "version": "latest",
        "downloadUrl": "",
        "versionsUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Terraform CLI to install. | string | latest | latest, 1.14.3, 1 |
| downloadUrl | The download URL to use for Terraform CLI binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/hashicorp-releases-generic-remote/terraform |
| versionsUrl | The URL to fetch the available Terraform CLI versions from. | string | &lt;empty&gt; |  |
