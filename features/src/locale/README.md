# Set a specific locale (locale)

Allows setting the locale.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/locale:0.1.0": {
        "locale": "en_US"
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| locale | The locale to set. | string | en_US | de_CH, en_US |

## Notes

Currently, you still need to add the appropriate language as environment variable to your container. You can do this in the `devcontainer.json` file like:
```json
{
    "containerEnv": {
        "LANG": "de_CH.UTF-8",
        "LANGUAGE": "de_CH.UTF-8",
        "LC_ALL": "de_CH.UTF-8"
    }
}
```

### System Compatibility

Debian, Ubuntu
