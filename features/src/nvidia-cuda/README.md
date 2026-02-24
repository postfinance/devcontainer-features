# NVIDIA CUDA (nvidia-cuda)

A package which installs NVIDIA CUDA.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/nvidia-cuda:0.1.0": {
        "version": "latest",
        "installLibraries": true,
        "installDevLibraries": true,
        "installCompiler": true,
        "installTools": true,
        "additionalCudaPackages": "",
        "downloadUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of NVIDIA CUDA to install. | string | latest | latest, 12.9, 13.1 |
| installLibraries | Installs all runtime CUDA Library packages. | boolean | true | true, false |
| installDevLibraries | Installs all development CUDA Library packages. | boolean | true | true, false |
| installCompiler | Installs all CUDA compiler packages. | boolean | true | true, false |
| installTools | Installs all CUDA command line and visual tools. | boolean | true | true, false |
| additionalCudaPackages | A comma-separated list of additional CUDA packages to install. Disable the other flags to have full control over what is installed. | string | &lt;empty&gt; |  |
| downloadUrl | The download URL to use for CUDA binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/nvidia-cuda-generic-remote/ |

## Notes

### System Compatibility

> Only amd64 is supported

Debian, Ubuntu 

### Accessed URLs

Needs access to the following URL for downloading and resolving:
* https://developer.download.nvidia.com

### Install the NVIDIA Container Toolkit

Follow [NVIDIA's instructions to install the NVIDIA Container Toolkit](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/overview.html) on your host machine. The NVIDIA Container Toolkit is available on a variety of Linux distributions. Make sure you have [installed the NVIDIA driver](https://docs.nvidia.com/datacenter/tesla/driver-installation-guide/index.html) for your Linux distribution before installing the NVIDIA Container Toolkit.

### Package Dependency Tree

Here is a list of the dependency tree of the packages in case you want to install only some of them.

cuda-toolkit
- cuda-compiler
  - cuda-cuobjdump
  - cuda-cuxxfilt
  - cuda-nvcc
  - cuda-nvprune
- cuda-libraries
  - libnpp
- cuda-libraries-dev
- cuda-tools
  - cuda-command-line-tools
    - cuda-nvtx
  - cuda-visual-tools
  - gds-tools
- cuda-documentation
- cuda-nvml-dev

### NVIDIA cuDNN

TODO: This is not yet supported

* libcudnn8 - cuDNN runtime libraries
* libcudnn8-dev - cuDNN development libraries and headers
* libcudnn9-cuda-11 - cuDNN runtime libraries for CUDA 11.8
* libcudnn9-dev-cuda-11 - cuDNN development headers and symlinks for CUDA 11.8
* libcudnn9-cuda-12 - cuDNN runtime libraries for CUDA 12.4
* libcudnn9-dev-cuda-12 - cuDNN development headers and symlinks for CUDA 12.4
