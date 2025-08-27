case $(uname -m | tr '[:upper:]' '[:lower:]') in
  x86_64*) ARCH=amd64 ;;
  arm*|aarch64*) ARCH=arm64 ;;
  *) echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

"./installer_$ARCH" \
    -version="${VERSION:-"lts"}" \
    -versionResolve="${VERSIONRESOLVE:-false}" \
    -npmVersion="${NPMVERSION:-"included"}" \
    -npmVersionResolve="${NPMVERSIONRESOLVE:-false}" \
    -yarnVersion="${YARNVERSION:-"none"}" \
    -yarnVersionResolve="${YARNVERSIONRESOLVE:-false}" \
    -pnpmVersion="${PNPMVERSION:-"none"}" \
    -pnpmVersionResolve="${PNPMVERSIONRESOLVE:-false}" \
    -downloadUrlBase="${DOWNLOADURLBASE:-""}" \
    -downloadUrlPath="${DOWNLOADURLPATH:-""}" \
    -versionsUrl="${VERSIONSURL:-""}" \
    -globalNpmRegistry="${GLOBALNPMREGISTRY:-""}"
