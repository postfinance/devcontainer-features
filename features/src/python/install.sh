. ./functions.sh

"./installer_$(detect_arch)" \
    -version="${VERSION:-"latest"}" \
    -downloadUrl="${DOWNLOADURL:-""}" \
    -pipIndex="${PIP_INDEX:-""}" \
    -pipIndexUrl="${PIP_INDEX_URL:-""}" \
    -pipTrustedHost="${PIP_TRUSTED_HOST:-""}"
