# CI Utility (ci-utility)

A package which installs various ci utility tools.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/ci-utility:0.1.0": {
        "yqVersion": "system-default",
        "gettextbaseVersion": "system-default",
        "yamllintVersion": "system-default",
        "sshpassVersion": "system-default"
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| yqVersion | The version of yq to install. | string | system-default | system-default, none |
| gettextbaseVersion | The version of gettext-base to install. | string | system-default | system-default, none |
| yamllintVersion | The version of yamllint to install. | string | system-default | system-default, none |
| sshpassVersion | The version of sshpass to install. | string | system-default | system-default, none |

## Notes

### System Compatibility

Debian (Bookworm, Trixie), Ubuntu
