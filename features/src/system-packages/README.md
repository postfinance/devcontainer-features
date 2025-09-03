# System Packages (system-packages)

Install arbitrary system packages using apt or apk.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/system-packages:0.1.0": {
        "packages": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| packages | Comma-separated list of system packages to install. | string | &lt;empty&gt; | curl,git,htop, sshpass,yamllint,yq,gettext-base |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine
