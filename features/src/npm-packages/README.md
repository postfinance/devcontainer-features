# NPM Packages (npm-packages)

Installs NPM packages globally.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/npm-packages:0.1.0": {
        "packages": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| packages | Comma-separated list of packages to install | string | &lt;empty&gt; | @devcontainers/cli,which@3.0.1 |
