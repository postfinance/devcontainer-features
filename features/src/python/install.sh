. ./functions.sh

"./installer_$(detect_arch)" \
    -version="${VERSION:-"latest"}" \
    -downloadUrl="${DOWNLOADURL:-""}" \
    -pipIndex="${PIP_INDEX:-""}" \
    -pipIndexUrl="${PIP_INDEX_URL:-""}" \
    -pipTrustedHost="${PIP_TRUSTED_HOST:-""}" \
    -uvVersion="${UV_VERSION:-"none"}" \
    -uvDownloadUrl="${UV_DOWNLOAD_URL:-""}" \
    -uvIndexUrl="${UV_INDEX_URL:-""}" \
    -uvExtraIndexUrl="${UV_EXTRA_INDEX_URL:-""}"
