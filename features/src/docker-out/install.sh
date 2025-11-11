. ./functions.sh

"./installer_$(detect_arch)" \
    -version="${VERSION:-"latest"}" \
    -composeVersion="${COMPOSEVERSION:-"latest"}" \
    -buildxVersion="${BUILDXVERSION:-"latest"}" \
    -configPath="${CONFIGPATH:-""}" \
    -downloadUrl="${DOWNLOADURL:-""}" \
    -versionsUrl="${VERSIONSURL:-""}" \
    -composeDownloadUrl="${COMPOSEDOWNLOADURL:-""}" \
    -buildxDownloadUrl="${BUILDXDOWNLOADURL:-""}"
