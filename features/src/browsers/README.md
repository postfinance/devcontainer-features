# Browsers (browsers)

A package which installs various browsers.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/browsers:1.1.1": {
        "chromeVersion": "none",
        "useChromeForTesting": true,
        "chromeDownloadUrl": "",
        "chromeVersionsUrl": "",
        "chromeTestingVersionsUrl": "",
        "firefoxVersion": "none",
        "firefoxDownloadUrl": "",
        "firefoxVersionsUrl": ""
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
