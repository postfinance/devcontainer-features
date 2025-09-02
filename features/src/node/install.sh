. ./functions.sh

"./installer_$(detect_arch)" \
    -version="${VERSION:-"lts"}" \
    -npmVersion="${NPMVERSION:-"included"}" \
    -yarnVersion="${YARNVERSION:-"none"}" \
    -pnpmVersion="${PNPMVERSION:-"none"}" \
    -corepackVersion="${COREPACKVERSION:-"none"}" \
    -downloadUrl="${DOWNLOADURL:-""}" \
    -versionsUrl="${VERSIONSURL:-""}" \
    -globalNpmRegistry="${GLOBALNPMREGISTRY:-""}"
