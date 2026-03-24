# .NET (dotnet)

A package which installs .NET SDKs, runtimes and workloads.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/dotnet:0.1.0": {
        "version": "lts",
        "additionalVersions": "",
        "dotnetRuntimeVersions": "",
        "aspNetCoreRuntimeVersions": "",
        "workloads": "",
        "downloadUrl": "",
        "versionsUrl": "",
        "nugetConfigPath": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | Select or enter a .NET SDK version. Use 'lts' for the latest LTS version, 'X.Y' or 'X.Y.Z' for a specific version. | string | lts | latest, lts, none, 8.0, 7.0, 6.0, 8.0.408 |
| additionalVersions | Enter additional .NET SDK versions, separated by commas. Use 'lts' for the latest LTS version, 'X.Y' or 'X.Y.Z' for a specific version. | string | &lt;empty&gt; | 7.0,8.0, 8.0.408 |
| dotnetRuntimeVersions | Enter additional .NET runtime versions, separated by commas. Use 'lts' for the latest LTS version, 'X.Y' or 'X.Y.Z' for a specific version. | string | &lt;empty&gt; | 8.0.15, 9.0, lts, 7.0 |
| aspNetCoreRuntimeVersions | Enter additional ASP.NET Core runtime versions, separated by commas. Use 'lts' for the latest LTS version, 'X.Y' or 'X.Y.Z' for a specific version. | string | &lt;empty&gt; | 8.0.15, lts, 7.0 |
| workloads | Enter additional .NET SDK workloads, separated by commas. Use 'dotnet workload search' to learn what workloads are available to install. | string | &lt;empty&gt; | wasm-tools, android, macos |
| downloadUrl | The download URL to use for Dotnet binaries. | string | &lt;empty&gt; |  |
| versionsUrl | The URL to use for fetching available Dotnet versions. | string | &lt;empty&gt; |  |
| nugetConfigPath | Path to a NuGet.Config file to copy into the container. This can be used to configure private package sources for the dotnet CLI. | string | &lt;empty&gt; |  |

## Customizations

### VS Code Extensions

- `ms-dotnettools.csharp`

## Dotnet Tools

If you need additional tools for example like the Powerapps CLI you can install them using `dotnet tool install --create-manifest-if-needed <Tool>`.
This installs the tool creates a manifest file: `.config/dotnet-tools.json`. 

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
