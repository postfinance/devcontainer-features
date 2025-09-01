. ./functions.sh

"./installer_$(detect_arch)" \
    -version="${VERSION:-"lts"}" \
    -npmVersion="${NPMVERSION:-"included"}" \
    -yarnVersion="${YARNVERSION:-"none"}" \
    -pnpmVersion="${PNPMVERSION:-"none"}" \
    -downloadUrl="${DOWNLOADURL:-""}" \
    -versionsUrl="${VERSIONSURL:-""}" \
    -globalNpmRegistry="${GLOBALNPMREGISTRY:-""}"
