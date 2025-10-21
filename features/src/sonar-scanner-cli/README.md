# SonarScanner CLI (sonar-scanner-cli)

A package which installs the SonarScanner CLI.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/sonar-scanner-cli:0.2.0": {
        "version": "latest",
        "includeJre": true,
        "downloadUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of SonarScanner CLI to install. | string | latest | latest, 6, 6.2.1.4610 |
| includeJre | A flag to indicate if the jre should be included in the download, otherwise a JRE needs to be installed in other ways. | boolean | true | true, false |
| downloadUrl | The download URL to use for the binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/sonarsource-generic-remote/Distribution/sonar-scanner-cli |

## Notes

### System Compatibility

Debian, Ubuntu

### Accessed Urls

Needs access to the following URL for downloading and resolving:
* https://binaries.sonarsource.com
* https://github.com
