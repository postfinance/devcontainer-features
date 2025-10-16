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
