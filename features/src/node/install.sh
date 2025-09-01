case $(uname -m | tr '[:upper:]' '[:lower:]') in
  x86_64*) ARCH=amd64 ;;
  arm*|aarch64*) ARCH=arm64 ;;
  *) echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

"./installer_$ARCH" \
    -version="${VERSION:-"lts"}" \
    -npmVersion="${NPMVERSION:-"included"}" \
    -yarnVersion="${YARNVERSION:-"none"}" \
    -pnpmVersion="${PNPMVERSION:-"none"}" \
    -downloadUrl="${DOWNLOADURL:-""}" \
    -versionsUrl="${VERSIONSURL:-""}" \
    -globalNpmRegistry="${GLOBALNPMREGISTRY:-""}"
