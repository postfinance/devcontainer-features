# Set a specific timezone (timezone)

Allows setting the timezone.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/timezone:0.1.0": {
        "timezone": "Etc/UTC"
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| timezone | The timezone to set. | string | Etc/UTC | Europe/Zurich, Etc/UTC |

## Notes

### System Compatibility

Debian, Ubuntu, Alpine
