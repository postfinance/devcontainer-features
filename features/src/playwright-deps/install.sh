. ./functions.sh

"./installer_$(detect_arch)" \
    -installChromiumDeps="${INSTALLCHROMIUMDEPS:-"true"}" \
    -installFirefoxDeps="${INSTALLFIREFOXDEPS:-"true"}" \
    -installWebkitDeps="${INSTALLWEBKITDEPS:-"true"}"
