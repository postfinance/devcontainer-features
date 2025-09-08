. ./functions.sh

"./installer_$(detect_arch)" \
  -version="${VERSION:-"latest"}" \
  -downloadUrl="${DOWNLOADURL:-""}" \
  -versionsUrl="${VERSIONSURL:-""}"
