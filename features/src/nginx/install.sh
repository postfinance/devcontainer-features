. ./functions.sh

"./installer_$(detect_arch)" \
  -version="${VERSION:-"latest"}" \
  -stableOnly="${STABLEONLY:-"false"}" \
  -downloadUrl="${DOWNLOADURL:-""}"
