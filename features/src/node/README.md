# Node.js (node)

A package which installs Node.js.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/node:0.3.0": {
        "version": "lts",
        "npmVersion": "included",
        "yarnVersion": "none",
        "pnpmVersion": "none",
        "corepackVersion": "none",
        "downloadUrl": "",
        "versionsUrl": "",
        "globalNpmRegistry": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Node.js to install. | string | lts | lts, latest, 20.11.1, 18.19.1 |
| npmVersion | The version of NPM to install. | string | included | included, latest, 10.5.0, 9.9.3 |
| yarnVersion | The version of Yarn to install. | string | none | none, latest, 1.22.22, 1.21.1 |
| pnpmVersion | The version of Pnpm to install. | string | none | none, latest, 9.14.2, 9 |
| corepackVersion | The version of corepack to install. | string | none | none, latest, 0.34.0, 0.29 |
| downloadUrl | The download URL to use for Node.js binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/nodejs-generic-remote/dist |
| versionsUrl | The URL to fetch the available Node.js versions from. | string | &lt;empty&gt; |  |
| globalNpmRegistry | The global NPM registry to use. | string | &lt;empty&gt; | https://mycompany.com/artifactory/api/npm/npm/ |

## Notes

### System Compatibility

Debian, Ubuntu

Alpine does not work as the binaries are compiled with glibc (instead of musl) which does not work on Alpine.
Some binaries could be taken from https://unofficial-builds.nodejs.org but ARM binaries are still missing.

### Accessed Urls

Needs access to the following URL for downloading and resolving:
* https://nodejs.org

### Corepack

:warning: Internet access is necessary for corepack to install your preferred package manager.

If you prefere to use internal sources, additional configuration is required. Add this to your `devcontainer.json`.

```json
  {
    "postCreateCommand": "no_proxy=.mycompany.com corepack install",
    "containerEnv": {
      "COREPACK_NPM_REGISTRY": "https://artifactory.mycompany.com/artifactory/api/npm/npm"
    }
  }
```

Notice the `no_proxy=.mycompany.com`; it is necessary because the package used by corepack does not follow the common rules for the `no_proxy` variable. See [Rob--W/proxy-from-env/issues#29](https://github.com/Rob--W/proxy-from-env/issues/29).

For **pnpm** to work with Artifactory, you have to additionally add this to the variables of your Dev Container:

```json
  {
    "containerEnv": {
      "COREPACK_INTEGRITY_KEYS": "0"
    }
  }
```

The reason for this are missing singatures in the Artifactory NPM API. See [nodejs/corepack#725](https://github.com/nodejs/corepack/issues/725)
