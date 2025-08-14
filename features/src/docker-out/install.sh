case $(uname -m | tr '[:upper:]' '[:lower:]') in
  x86_64*) ARCH=amd64 ;;
  arm*|aarch64*) ARCH=arm64 ;;
  *) echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

"./installer_$ARCH" \
    -version="${VERSION:-"latest"}" \
    -versionResolve="${VERSIONRESOLVE:-"false"}" \
    -composeVersion="${COMPOSEVERSION:-"latest"}" \
    -composeVersionResolve="${COMPOSEVERSIONRESOLVE:-"false"}" \
    -buildxVersion="${BUILDXVERSION:-"latest"}" \
    -buildxVersionResolve="${BUILDXVERSIONRESOLVE:-"false"}" \
    -downloadUrlBase="${DOWNLOADURLBASE:-""}" \
    -downloadUrlPath="${DOWNLOADURLPATH:-""}" \
    -versionsUrl="${VERSIONSURL:-""}" \
    -composeDownloadUrlBase="${COMPOSEDOWNLOADURLBASE:-""}" \
    -composeDownloadUrlPath="${COMPOSEDOWNLOADURLPATH:-""}" \
    -buildxDownloadUrlBase="${BUILDXDOWNLOADURLBASE:-""}" \
    -buildxDownloadUrlPath="${BUILDXDOWNLOADURLPATH:-""}"
