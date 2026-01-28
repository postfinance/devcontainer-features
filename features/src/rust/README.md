# Rust (rust)

A package which installs Rust, common Rust utilities and their required dependencies.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/rust:0.1.0": {
        "version": "latest",
        "rustupVersion": "latest",
        "profile": "minimal",
        "components": "rustfmt,rust-analyzer,rust-src,clippy",
        "enableWindowsTarget": false
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Rust to install. | string | latest | latest, 1.93.0 |
| rustupVersion | The version of rustup to install. | string | latest | latest, 1.27.1 |
| profile | The rustup profile to install. | string | minimal | minimal, default, complete |
| components | A comma separated list with components that should be installed. | string | rustfmt,rust-analyzer,rust-src,clippy | , rustfmt,rust-analyzer, rls,rust-analysis |
| enableWindowsTarget | A flag to indicate if the Windows target (and needed tools) should be installed. | boolean | false | true, false |

## Customizations

### VS Code Extensions

- `vadimcn.vscode-lldb`
- `rust-lang.rust-analyzer`
- `tamasfe.even-better-toml`
- `serayuzgur.crates`

## Notes

### System Compatibility

Debian, Ubuntu

### Accessed Urls

Needs access to the following URL for downloading and resolving:
* https://static.rust-lang.org
