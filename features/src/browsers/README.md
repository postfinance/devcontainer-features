# Browsers (browsers)

Installs various browsers and their dependencies.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/browsers:0.1.0": {
        "chromeVersion": "none",
        "useChromeForTesting": true,
        "chromeDownloadUrl": "",
        "chromeVersionsUrl": "",
        "chromeTestingVersionsUrl": "",
        "firefoxVersion": "none",
        "firefoxDownloadUrl": "",
        "firefoxVersionsUrl": "",
        "firefoxVersionResolve": false
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| chromeVersion | The version of the Chrome to install. | string | none | none, latest, 126 |
| useChromeForTesting | A flag to indicate if the Chrome for Testing or the default chrome should be used. | boolean | true | true, false |
| chromeDownloadUrl | Override Chrome download base URL. | string | &lt;empty&gt; |  |
| chromeVersionsUrl | Override Chrome versions URL. | string | &lt;empty&gt; |  |
| chromeTestingVersionsUrl | Override Chrome for Testing versions URL. | string | &lt;empty&gt; |  |
| firefoxVersion | The version of the Firefox to install. | string | none | none, latest, 128 |
| firefoxDownloadUrl | Override Firefox download base URL. | string | &lt;empty&gt; |  |
| firefoxVersionsUrl | Override Firefox versions URL. | string | &lt;empty&gt; |  |
| firefoxVersionResolve | If true, resolves partial Firefox versions (e.g. 142.0) to the highest available patch version (e.g. 142.0.3). | boolean | false | true, false |

## Notes

### System Compatibility

Debian, Ubuntu

### Accessed Urls

Needs access to the following URL for downloading and resolving:
* https://dl.google.com
* https://versionhistory.googleapis.com
* https://googlechromelabs.github.io
* https://download-installer.cdn.mozilla.net
* https://product-details.mozilla.org
