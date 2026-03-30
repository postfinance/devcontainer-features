## Dotnet Tools

If you need additional tools, such as the PowerApps CLI, you can install them by running `dotnet tool install --create-manifest-if-needed <Tool>`.
This installs the tool and creates a manifest file: `.config/dotnet-tools.json`. 

```json
{
  "version": 1,
  "isRoot": true,
  "tools": {
    "microsoft.powerapps.cli.tool": {
      "version": "1.43.6",
      "commands": [
        "pac"
      ],
      "rollForward": false
    }
  }
}
```

After this step, the tool can be invoked using `dotnet <command>`.

If you already have a manifest, all tools can be installed using `dotnet tool restore`.

To do that automatically, include the command in your `devcontainer.json` like this:
```json
"postCreateCommand": "dotnet tool restore"
```

### System Compatibility

Debian, Ubuntu, Alpine
