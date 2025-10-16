# Playwright Dependencies (playwright-deps)

A package which installs the needed dependencies to run Playwright.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/playwright-deps:0.1.0": {
        "installChromiumDeps": true,
        "installFirefoxDeps": true,
        "installWebkitDeps": true
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| installChromiumDeps | A flag to indicate if the dependencies for Chromium should be installed. | boolean | true | true, false |
| installFirefoxDeps | A flag to indicate if the dependencies for Firefox should be installed. | boolean | true | true, false |
| installWebkitDeps | A flag to indicate if the dependencies for Webkit should be installed. | boolean | true | true, false |

## Notes

### System Compatibility

Debian, Ubuntu

### Playwright Installation

Playwright itself is not installed with this feature. This should be managed and installed with your package manager (`npm`, `yarn`, ...).

### Browser Paths

By default, the browsers are in a folder that is cleared when the dev-container is rebuilt.

In order to persist them, you can set the path to the browsers to a folder in the workspace.

To do this, add the following environment variable to your container:
```
"containerEnv": {
  "PLAYWRIGHT_BROWSERS_PATH": "${containerWorkspaceFolder}/.ms-playwright"
}
```

Also don't forget to add `.ms-playwright` to your `.gitignore` file.

### Download URL from within Playwright

To change the downloads that Playwright does to a custom remote, you can adjust the `PLAYWRIGHT_DOWNLOAD_HOST` environment variable to eg. `https://mycompany.com/artifactory/playwright-remote/`. You might also want to increase the `PLAYWRIGHT_DOWNLOAD_CONNECTION_TIMEOUT` to something like `300000`.
