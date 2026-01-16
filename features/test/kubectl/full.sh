#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_version "$(kubectl version --client | head -1)" "Client Version: v1.31.2"
check_version "$(kubectx --version)" "0.9.5"
check_version "$(kubens --version)" "0.9.5"
check_version "$(k9s version --short | head -1)" "v0.32.5"
check_version "$(helm version --short)" "v3.16.2+g13654a5"
check_version "$(kustomize version)" "v5.4.2"
check_version "$(kubeconform -v)" "v0.6.6"
check_version "$(kube-score version)" "kube-score version: 1.18.0, commit: 0fb5f668e153c22696aa75ec769b080c41b5dd3d, built: 2024-02-05T14:13:15Z"
