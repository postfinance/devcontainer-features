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
