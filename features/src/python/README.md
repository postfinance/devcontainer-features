# Python (python)

A package which installs Python.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/python:0.2.0": {
        "version": "latest",
        "downloadUrl": "",
        "pipIndex": "",
        "pipIndexUrl": "",
        "pipTrustedHost": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Python to install. | string | latest | latest, 3.12, 3.9.19 |
| downloadUrl | The download URL to use. | string | &lt;empty&gt; | https://mycompany.com/artifactory/python-generic-ftp-remote, https://mycompany.com/artifactory/python-generic-remote/ftp/python |
| pipIndex | The pip index to use (used by search). | string | &lt;empty&gt; | https://mycompany.com/artifactory/api/pypi/python/simple, https://mycompany.com/nexus/repository/pypi-group/pypi |
| pipIndexUrl | The pip index URL to use (used by install). | string | &lt;empty&gt; | https://mycompany.com/artifactory/api/pypi/python/simple, https://mycompany.com/nexus/repository/pypi-group/simple |
| pipTrustedHost | The pip trusted host to use. | string | &lt;empty&gt; | mycompany.com, artifactory.mycompany.com, nexus.mycompany.com |

## Customizations

### VS Code Extensions

- `ms-python.python`
- `ms-python.vscode-pylance`

## Notes

### .0 versions for old Python versions

Python <= 3.2 didn't use correct semver and therefore the .0 version was missing the .0 (eg. 3.1.0 was only released as 3.1).

Because of this, it is not possible to install an exact .0 version for those Python releases, it will take the newest patch release in that case.

Example:
* version 3.1 => would install 3.1.5
* version 3.1.0 => would unfortunately fail
