# kubectl (kubectl)

Installs kubectl and other tools for managing kubernetes.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/kubectl:0.1.0": {
        "version": "latest",
        "kubectxVersion": "latest",
        "kubensVersion": "latest",
        "k9sVersion": "none",
        "helmVersion": "none",
        "kustomizeVersion": "none",
        "kubeconformVersion": "none",
        "kubescoreVersion": "none",
        "downloadUrl": "",
        "kubectxDownloadUrl": "",
        "kubensDownloadUrl": "",
        "k9sDownloadUrl": "",
        "helmDownloadUrl": "",
        "kustomizeDownloadUrl": "",
        "kubeconformDownloadUrl": "",
        "kubescoreDownloadUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of kubectl to install. | string | latest | latest, none, 1.30.0 |
| kubectxVersion | The version of kubectx to install. | string | latest | latest, none, 0.9.5 |
| kubensVersion | The version of kubens to install. | string | latest | latest, none, 0.9.5 |
| k9sVersion | The version of k9s to install. | string | none | latest, none, 0.32.5 |
| helmVersion | The version of helm to install. | string | none | latest, none, 3.16.2 |
| kustomizeVersion | The version of kustomize to install. | string | none | latest, none, 5.4.2 |
| kubeconformVersion | The version of kubeconform to install. | string | none | latest, none, 0.6.6 |
| kubescoreVersion | The version of kube-score to install. | string | none | latest, none, 1.18.0 |
| downloadUrl | The download URL to use for kubectl binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/dlk8sio-generic-remote/release |
| kubectxDownloadUrl | The download URL to use for kubectx binaries. | string | &lt;empty&gt; |  |
| kubensDownloadUrl | The download URL to use for kubens binaries. | string | &lt;empty&gt; |  |
| k9sDownloadUrl | The download URL to use for k9s binaries. | string | &lt;empty&gt; |  |
| helmDownloadUrl | The download URL to use for helm binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/gethelmsh-generic-remote |
| kustomizeDownloadUrl | The download URL to use for kustomize binaries. | string | &lt;empty&gt; |  |
| kubeconformDownloadUrl | The download URL to use for kubeconform binaries. | string | &lt;empty&gt; |  |
| kubescoreDownloadUrl | The download URL to use for kube-score binaries. | string | &lt;empty&gt; |  |
