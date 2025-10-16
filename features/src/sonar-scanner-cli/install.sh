. ./functions.sh

"./installer_$(detect_arch)" \
    -version="${VERSION:-"latest"}" \
    -includeJre="${INCLUDEJRE:-"true"}"
